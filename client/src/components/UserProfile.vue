<template>
  <!-- Profile page -->
  <div v-if="!isProfileEditMode" class="user-profile">
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
        <span class="profile-edit" @click="switchToEditMode">Изм.</span>
      </div>

      <!-- Profile content -->
      <div class="profile-content">
        <div class="profile-main-info">
          <div class="avatar">
            <img
              v-if="user?.avatarUrl"
              :src="user.avatarUrl"
              alt="Avatar"
              class="avatar-image"
              @error="onAvatarError"
            />
            <span v-else class="avatar-initials">{{ userInitials }}</span>
          </div>
          <div class="user-name-details">
            <div class="user-full-name" v-if="userFullName">
              {{ userFullName }}
            </div>
            <div class="username">{{ user?.username }}</div>
          </div>
        </div>

        <div class="profile-bio">Здесь будет биография пользователя.</div>
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
          <div class="menu-item-container">
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
          <div class="menu-item-container">
            <button class="profile-list-button">
              <span>Помощь</span>
            </button>
            <button class="profile-list-button">
              <span>О приложении</span>
            </button>
          </div>

          <!-- Logout -->
          <div class="menu-item-container">
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

  <!-- Edit profile page -->
  <UserProfileEdit
    v-else
    :current-user="user"
    @cancel="switchToViewMode"
    @save="handleSave"
  />
</template>

<script>
import { mapGetters, mapActions } from "vuex";
import UserProfileEdit from "./UserProfileEdit.vue";
export default {
  name: "UserProfile",
  components: { UserProfileEdit },
  data() {
    return {
      isProfileEditMode: false,
    };
  },
  computed: {
    ...mapGetters("auth", ["user"]),
    ...mapGetters("theme", ["isDarkTheme"]),

    userFullName() {
      if (!this.user) return "";
      const { firstName, lastName } = this.user;
      return [firstName, lastName].filter(Boolean).join(" ").trim() || null;
    },
    userInitials() {
      // Используем this.user, полученный из mapGetters
      const first = this.user?.firstName?.[0] || "";
      const last = this.user?.lastName?.[0] || "";
      const usernameInitial = this.user?.username?.[0] || "?";
      const fallbackInitial =
        usernameInitial === "?" ? "?" : usernameInitial.toUpperCase();
      // Собираем инициалы: первые буквы имени и фамилии,
      // или первая буква имени пользователя, если имя/фамилия не заданы
      const initials = (first + last).toUpperCase();
      return initials || fallbackInitial;
    },
  },
  watch: {
    user(newUser) {
      if (!newUser) {
        this.isProfileEditMode = false;
      }
    },
  },
  methods: {
    ...mapActions("theme", ["toggleTheme"]),

    onAvatarError(event) {
      console.warn("UserProfile: Failed to load avatar image.", event);
    },
    
    closeProfile() {
      this.$emit("close");
    },
    switchToEditMode() {
      this.isProfileEditMode = true;
    },
    switchToViewMode() {
      this.isProfileEditMode = false;
    },
    handleSave() {
      this.isProfileEditMode = false;
    },
    addAccount() {
      // add account action here
    },
    openSettings() {
      // open settings page
    },
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
