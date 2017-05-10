package captiveportal::PacketFence::Controller::LostStolen;
use Moose;
use namespace::autoclean;
use pf::violation;

BEGIN { extends 'captiveportal::Base::Controller'; }

=head1 NAME

captiveportal::PacketFence::Controller::LostStolen - Catalyst Controller

=head1 DESCRIPTION

Catalyst Controller.

=head1 METHODS

=cut

=head2 index

=cut

sub index : Path : Args(1) {
    my ( $self, $c, $mac ) = @_;
    my $portalSession = $c->portalSession;

    my $violation = violation_trigger($mac);

    if ($violation) {

        # There is  to trigger the violation Lost or Stolen
        # Put the device in isolation and alerts user + admin
        my $vid   = $violation->{'vid'};
        my $class = class_view($vid);
        $c->stash(
            violation_id => $vid,
            enable_text  => $class->{button_text},
            template     => 'lost_stolen.html',
        );
    } else {
        $self->showError( $c, "error to apply the status: Lost or Stolen" );
    }
}

=head1 AUTHOR

Inverse inc. <info@inverse.ca>

=head1 COPYRIGHT

Copyright (C) 2005-2017 Inverse inc.

=head1 LICENSE

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301,
USA.

=cut

1;
