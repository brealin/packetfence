<template>
  <b-button :style="{ minWidth: btnWidth }" ref="saveButton" :variant="variant" type="submit" v-bind="$attrs" v-on="forwardListeners">
    <icon name="circle-notch" spin v-if="isLoading"></icon>
    <template v-else>
      <slot>{{ $t('Save') }}</slot>
    </template>
  </b-button>
</template>

<script>
export default {
  name: 'pf-button-save',
  props: {
    isLoading: {
      type: Boolean,
      default: false
    },
    variant: {
      type: String,
      default: 'primary'
    }
  },
  data () {
    return {
      btnWidth: 0
    }
  },
  computed: {
    forwardListeners () {
      const { input, ...listeners } = this.$listeners
      return listeners
    }
  },
  watch: {
    isLoading: {
      handler: function (newValue) {
        if (newValue) {
          this.btnWidth = (this.$refs.saveButton.clientWidth + 2) + 'px'
        }
      }
    }
  }
}
</script>
