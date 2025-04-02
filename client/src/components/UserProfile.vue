<template>
  <!-- Profile page -->
  <div class="user-profile">
    <div class="content">
      <!-- Profile header -->
      <div class="profile-header">
        <button class="close-button" @click="closeProfile">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="28"
            height="28"
            fill="none"
            viewBox="0 0 24 24"
          >
            <line
              x1="18"
              y1="6"
              x2="6"
              y2="18"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
            />
            <line
              x1="6"
              y1="6"
              x2="18"
              y2="18"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
            />
          </svg>
        </button>

        <!-- Edit profile -->
        <span class="profile-edit">Изм.</span>
      </div>

      <!-- Profile content -->
      <div class="profile-content">
        <div class="avatar"></div>
        <div class="username">@{{ user?.username }}</div>
        <div class="user-email">{{ user?.email }}</div>
      </div>

      <!-- Profile list preferences -->
      <div class="profile-list-scroll">
        <div class="profile-list">
          <!-- Add account -->
          <div class="menu-item-container">
            <button class="profile-list-button" @click="addAccount">
              <svg
                class="plus-user-icon"
                xmlns="http://www.w3.org/2000/svg"
                width="26"
                height="26"
                fill="none"
                viewBox="0 0 24 24"
              >
                <line
                  x1="12"
                  y1="5"
                  x2="12"
                  y2="19"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                />
                <line
                  x1="5"
                  y1="12"
                  x2="19"
                  y2="12"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                />
              </svg>

              <span>Добавить аккаунт</span>
            </button>
          </div>

          <!-- Profile and app settings -->
          <div class="menu-item-container has-divider">
            <button class="profile-list-button" @click="toggleTheme">
              <span>Внешний вид</span>
            </button>
            <button class="profile-list-button">
              <span>Уведомления</span>
            </button>
            <button class="profile-list-button">
              <span>Конфиденциальность</span>
            </button>
            <button class="profile-list-button" @click="openSettings">
              <span>Настройки</span>
            </button>
          </div>

          <!-- Info -->
          <div class="menu-item-container has-divider">
            <button class="profile-list-button">
              <span>Помощь</span>
            </button>
            <button class="profile-list-button">
              <span>О приложении</span>
            </button>
          </div>

          <!-- Logout -->
          <div class="menu-item-container has-divider">
            <button class="profile-list-button logout-button" @click="logout">
              <svg class="menu-icon" viewBox="0 0 24 24" fill="currentColor">
                <path
                  d="M17 8l-1.41 1.41L17.17 11H9v2h8.17l-1.58 1.58L17 16l4-4-4-4zM5 5h7V3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h7v-2H5V5z"
                />
              </svg>
              <span>Выход</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { mapGetters, mapActions } from "vuex";
export default {
  name: "UserProfile",
  props: {
    visible: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    ...mapGetters("auth", ["user"]),
    ...mapGetters("theme", ["isDarkTheme"]),
  },
  methods: {
    ...mapActions("theme", ["toggleTheme"]),

    closeProfile() {
      this.$emit("close");
    },
    addAccount() {
      // add account action here
    },
    openSettings() {
      // open settings page
    },
    // logout
    logout() {
      this.$store.dispatch("auth/logout");
      this.$router.push({ name: "AuthPage" });
    },
  },
};
</script>

<style scoped>
@import "@/styles/userprofile.css";
</style>
