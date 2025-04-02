<template>
  <div class="username-tooltip">
    <!-- Info button -->
    <button type="button" class="info-icon" @click="toggleUsernameInfo">
      <svg class="info-svg" viewBox="0 0 24 24" fill="none">
        <circle
          cx="12"
          cy="12"
          r="8.5"
          stroke="currentColor"
          stroke-width="1.5"
        />
        <path
          d="M12 11V16"
          stroke="currentColor"
          stroke-width="1.5"
          stroke-linecap="round"
        />
        <circle cx="12" cy="8" r="1" fill="currentColor" />
      </svg>
    </button>

    <!-- Username tooltip -->
    <transition name="tooltip">
      <div
        v-if="showUsernameInfo"
        class="info-tooltip"
        :class="{ dark: isDarkTheme }"
        ref="tooltip"
      >
        <div class="tooltip-content">
          С Вами можно будет связаться по выбранному имени пользователя.
          <br />
          Имя пользователя может содержать:
          <ul>
            <li>Только латинские буквы</li>
            <li>Цифры и нижние подчёркивания</li>
            <li>Не менее 5 символов</li>
          </ul>
        </div>
      </div>
    </transition>
  </div>
</template>

<script>
export default {
  name: "UsernameTooltip",
  props: {
    isDarkTheme: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      showUsernameInfo: false,
    };
  },
  mounted() {
    document.addEventListener("click", this.handleClickOutside);
  },
  beforeUnmount() {
    document.removeEventListener("click", this.handleClickOutside);
  },
  methods: {
    toggleUsernameInfo() {
      this.showUsernameInfo = !this.showUsernameInfo;
      if (this.showUsernameInfo) {
        this.$nextTick(() => {
          this.updateTooltipPosition();
        });
      }
    },
    handleClickOutside(e) {
      if (this.showUsernameInfo) {
        const tooltipEl = this.$refs.tooltip;
        const iconEl = this.$el.querySelector(".info-icon");
        if (
          tooltipEl &&
          !tooltipEl.contains(e.target) &&
          iconEl &&
          !iconEl.contains(e.target)
        ) {
          this.showUsernameInfo = false;
        }
      }
    },
    updateTooltipPosition() {
      const iconEl = this.$el.querySelector(".info-icon");
      if (!iconEl) return;
      const rect = iconEl.getBoundingClientRect();
      const left = rect.right + 8;
      const top = rect.top;
      this.tooltipStyle = {
        position: "fixed",
        top: `${top}px`,
        left: `${left}px`,
      };
    },
  },
};
</script>

<style scoped>
@import "@/styles/authpage.css";
</style>
