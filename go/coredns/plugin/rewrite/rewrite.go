package rewrite

import (
	"fmt"
	"strings"

	"github.com/inverse-inc/packetfence/go/coredns/plugin"

	"github.com/miekg/dns"

	"golang.org/x/net/context"
)

// Result is the result of a rewrite
type Result int

const (
	// RewriteIgnored is returned when rewrite is not done on request.
	RewriteIgnored Result = iota
	// RewriteDone is returned when rewrite is done on request.
	RewriteDone
	// RewriteStatus is returned when rewrite is not needed and status code should be set
	// for the request.
	RewriteStatus
)

// These are defined processing mode.
const (
	// Processing should stop after completing this rule
	Stop = "stop"
	// Processing should continue to next rule
	Continue = "continue"
)

// Rewrite is plugin to rewrite requests internally before being handled.
type Rewrite struct {
	Next     plugin.Handler
	Rules    []Rule
	noRevert bool
}

// ServeDNS implements the plugin.Handler interface.
func (rw Rewrite) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	wr := NewResponseReverter(w, r)
	for _, rule := range rw.Rules {
		switch result := rule.Rewrite(w, r); result {
		case RewriteDone:
			if rule.Mode() == Stop {
				if rw.noRevert {
					return plugin.NextOrFailure(rw.Name(), rw.Next, ctx, w, r)
				}
				return plugin.NextOrFailure(rw.Name(), rw.Next, ctx, wr, r)
			}
		case RewriteIgnored:
			break
		case RewriteStatus:
			// only valid for complex rules.
			// if cRule, ok := rule.(*ComplexRule); ok && cRule.Status != 0 {
			// return cRule.Status, nil
			// }
		}
	}
	return plugin.NextOrFailure(rw.Name(), rw.Next, ctx, w, r)
}

// Name implements the Handler interface.
func (rw Rewrite) Name() string { return "rewrite" }

// Rule describes a rewrite rule.
type Rule interface {
	// Rewrite rewrites the current request.
	Rewrite(dns.ResponseWriter, *dns.Msg) Result
	// Mode returns the processing mode stop or continue
	Mode() string
}

func newRule(args ...string) (Rule, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("no rule type specified for rewrite")
	}

	arg0 := strings.ToLower(args[0])
	var ruleType string
	var expectNumArgs, startArg int
	mode := Stop
	switch arg0 {
	case Continue:
		mode = arg0
		ruleType = strings.ToLower(args[1])
		expectNumArgs = len(args) - 1
		startArg = 2
	case Stop:
		ruleType = strings.ToLower(args[1])
		expectNumArgs = len(args) - 1
		startArg = 2
	default:
		// for backward compatibility
		ruleType = arg0
		expectNumArgs = len(args)
		startArg = 1
	}

	if ruleType != "edns0" && expectNumArgs != 3 {
		return nil, fmt.Errorf("%s rules must have exactly two arguments", ruleType)
	}
	switch ruleType {
	case "name":
		return newNameRule(args[startArg], args[startArg+1])
	case "class":
		return newClassRule(args[startArg], args[startArg+1])
	case "type":
		return newTypeRule(args[startArg], args[startArg+1])
	case "edns0":
		return newEdns0Rule(mode, args[startArg:]...)
	default:
		return nil, fmt.Errorf("invalid rule type %q", args[0])
	}
}
