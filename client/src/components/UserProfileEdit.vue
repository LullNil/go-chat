<template>
  <!-- Edit Profile page -->
  <div class="user-profile">
    <div class="content">
      <!-- Profile header -->
      <div class="profile-header">
        <!-- Cancel button -->
        <button
          class="header-button cancel-button"
          @click="cancelEdit"
          :disabled="isSaving"
        >
          Отмена
        </button>

        <!-- Edit Profile title -->
        <div class="header-title-wrapper">
          <span class="header-title">Редактирование</span>
        </div>

        <!-- Save button -->
        <button
          class="header-button save-button"
          @click="saveProfile"
          :disabled="isSaving"
          :class="{ 'is-loading': isSaving }"
        >
          <span v-if="!isSaving">Готово</span>
          <!-- Loading indicator -->
          <svg v-else class="button-spinner-svg" viewBox="25 25 50 50">
            <circle
              class="spinner-path"
              cx="50"
              cy="50"
              r="20"
              fill="none"
              stroke-width="6"
              stroke-miterlimit="10"
            />
          </svg>
        </button>
      </div>

      <!-- Profile list -->
      <div class="profile-list-scroll edit-scroll">
        <form @submit.prevent="saveProfile" class="edit-profile-form">
          <fieldset
            :disabled="isSaving"
            class="edit-form-fieldset"
            ref="formContainer"
          >
            <div class="edit-profile-top-section">
              <!-- Avatar -->
              <div class="edit-avatar-wrapper">
                <div
                  class="avatar large-avatar"
                  @click="changeAvatar"
                  :class="{ disabled: isSaving }"
                  role="button"
                  tabindex="0"
                  @keydown.enter="changeAvatar"
                  @keydown.space="changeAvatar"
                >
                  <img
                    v-if="currentUser?.avatarUrl"
                    :src="currentUser.avatarUrl"
                    alt="Avatar"
                    class="avatar-image"
                    @error="onAvatarError"
                  />
                  <span v-else class="avatar-initials">{{ userInitials }}</span>

                  <div class="avatar-overlay">
                    <svg
                      class="edit-icon"
                      xmlns="http://www.w3.org/2000/svg"
                      viewBox="0 0 24 24"
                      fill="currentColor"
                      width="24"
                      height="24"
                    >
                      <path d="M0 0h24v24H0V0z" fill="none" />
                      <path
                        d="M3 17.25V21h3.75L17.81 9.94l-3.75-3.75L3 17.25zM20.71 7.04c.39-.39.39-1.02 0-1.41l-2.34-2.34c-.39-.39-1.02-.39-1.41 0l-1.83 1.83 3.75 3.75 1.83-1.83z"
                      />
                    </svg>
                    <span class="overlay-text">Изменить фото</span>
                  </div>
                </div>
                <input
                  type="file"
                  ref="avatarInput"
                  @change="handleAvatarChange"
                  accept="image/png, image/jpeg, image/webp"
                  style="display: none"
                />
              </div>

              <div class="edit-name-inputs menu-item-container edit-form-group">
                <div class="input-wrapper" @click="focusInput('firstName')">
                  <label for="firstName" class="input-label-overlay"></label>
                  <input
                    type="text"
                    id="firstName"
                    placeholder="Имя"
                    v-model.trim="editableUser.firstName"
                    autocomplete="given-name"
                  />
                </div>
                <div class="input-wrapper" @click="focusInput('lastName')">
                  <label for="lastName" class="input-label-overlay"></label>
                  <input
                    type="text"
                    id="lastName"
                    placeholder="Фамилия"
                    v-model.trim="editableUser.lastName"
                    autocomplete="family-name"
                  />
                </div>
              </div>
            </div>

            <div class="edit-form-list">
              <div class="form-section">
                <div class="menu-item-container edit-form-group">
                  <div class="input-wrapper" @click="focusInput('username')">
                    <label for="username" class="input-label-overlay"></label>
                    <input
                      type="text"
                      id="username"
                      placeholder="Имя пользователя"
                      v-model.trim="editableUser.username"
                      autocomplete="username"
                    />
                  </div>
                  <div class="input-wrapper" @click="focusInput('email')">
                    <label for="email" class="input-label-overlay"></label>
                    <input
                      type="email"
                      id="email"
                      placeholder="Email"
                      v-model.trim="editableUser.email"
                      autocomplete="email"
                    />
                  </div>
                </div>
                <p class="form-group-info">
                  Имя пользователя должно быть уникальным.
                </p>
              </div>

              <div class="form-section">
                <div class="menu-item-container edit-form-group">
                  <div
                    class="input-wrapper textarea-wrapper"
                    @click="focusInput('bio')"
                  >
                    <label for="bio" class="input-label-overlay"></label>
                    <textarea
                      id="bio"
                      ref="bioTextarea"
                      placeholder="О себе (необязательно)"
                      v-model="editableUser.bio"
                      @input="adjustTextareaHeight"
                    ></textarea>
                  </div>
                </div>
                <p class="form-group-info">
                  Краткое описание, которое будет видно в вашем профиле.
                </p>
              </div>
            </div>
          </fieldset>
        </form>
      </div>
    </div>

    <!-- Error message -->
    <div v-if="errorMessage" class="profile-error-message">
      {{ errorMessage }}
    </div>
  </div>
</template>

<script>
import { updateUserProfile } from "@/composables/useAuthAPI";
import { mapActions } from "vuex";

/**
 * Compares relevant fields between original and current user data objects.
 * @param {object} original - The original user data object.
 * @param {object} current - The current editable user data object.
 * @returns {boolean} True if changes are detected, false otherwise.
 */
function checkProfileChanges(original, current) {
  const fields = ["firstName", "lastName", "username", "email", "bio"];
  for (const field of fields) {
    const currentValue = current[field] ?? "";
    const originalValue = original[field] ?? "";
    if (currentValue !== originalValue) {
      return true; // Change detected
    }
  }
  return false; // No changes detected
}

export default {
  name: "UserProfileEdit",
  props: {
    /**
     * The current user object (with camelCase keys) from the Vuex store.
     */
    currentUser: {
      type: Object,
      required: true,
    },
  },
  data() {
    // Initialize editable and original data based on the currentUser prop
    const userData = {
      firstName: this.currentUser?.firstName || "",
      lastName: this.currentUser?.lastName || "",
      username: this.currentUser?.username || "",
      email: this.currentUser?.email || "",
      bio: this.currentUser?.bio || "",
    };
    return {
      // Reactive copy for user edits
      editableUser: { ...userData },
      // Static copy for comparison, do not modify!
      originalUser: { ...userData },
      isSaving: false,
      errorMessage: null,
      // TODO: Add state for the selected avatar file to upload
      // selectedAvatarFile: null,
    };
  },
  computed: {
    /**
     * Generates user initials for the avatar placeholder.
     */
    userInitials() {
      const first = this.currentUser?.firstName?.[0] || "";
      const last = this.currentUser?.lastName?.[0] || "";
      const usernameInitial = this.currentUser?.username?.[0] || "?";
      // Use cyrillic '?' if appropriate for the locale, otherwise standard '?'
      const fallbackInitial =
        usernameInitial === "?" ? "?" : usernameInitial.toUpperCase();
      return (first + last).toUpperCase() || fallbackInitial;
    },
  },
  watch: {
    // Watch for changes in the bio to adjust textarea height
    "editableUser.bio"() {
      this.$nextTick(() => {
        this.adjustTextareaHeight();
      });
    },
  },
  mounted() {
    // Adjust textarea height on initial mount if there's content
    this.$nextTick(() => {
      this.adjustTextareaHeight();
    });
  },
  methods: {
    // Map Vuex action for updating user data in the store
    ...mapActions("auth", ["updateUser"]),

    /**
     * Emits cancel event if not currently saving.
     */
    cancelEdit() {
      if (this.isSaving) return;
      this.$emit("cancel");
    },

    /**
     * Handles the profile saving process.
     * Checks for changes, calls the API, updates the store, and handles errors.
     */
    async saveProfile() {
      if (this.isSaving) return;

      // Check for changes before attempting to save
      const hasTextChanges = checkProfileChanges(
        this.originalUser,
        this.editableUser
      );
      // TODO: Add check for selectedAvatarFile if implementing avatar upload
      const hasAvatarChanged = false; // Placeholder

      if (!hasTextChanges && !hasAvatarChanged) {
        this.$emit("cancel");
        return;
      }

      this.isSaving = true;
      this.errorMessage = null;

      // Prepare data payload (snake_case for API)
      const profileDataToUpdate = {
        first_name: this.editableUser.firstName,
        last_name: this.editableUser.lastName,
        username: this.editableUser.username,
        email: this.editableUser.email,
        bio: this.editableUser.bio,
        // TODO: Handle avatar update
      };

      try {
        // Call API to update profile data
        const updatedUserDataFromApi = await updateUserProfile(
          profileDataToUpdate
        );

        // Update Vuex Store
        this.updateUser(updatedUserDataFromApi);

        // Emit save event to parent
        this.$emit("save", updatedUserDataFromApi);
      } catch (error) {
        // Handle API or other errors
        console.error("UserProfileEdit: Error updating profile:", error);
        this.errorMessage = this.formatErrorMessage(error);
      } finally {
        // Reset saving state
        this.isSaving = false;
      }
    },

    /**
     * Formats error messages from various possible error structures.
     * @param {Error|Object} error - The error object caught.
     * @returns {string} A user-friendly error message.
     */
    formatErrorMessage(error) {
      // Use default message
      let message = "Не удалось сохранить изменения. Попробуйте снова.";
      if (error?.data?.message) {
        message = error.data.message; // Use server message if available
      } else if (typeof error?.message === "string") {
        try {
          const errData = JSON.parse(error.message);
          message = errData.message || message;
        } catch (e) {
          /* ignore parsing error */
        }
      } else if (error?.statusText) {
        message = `Ошибка ${error.status || ""}: ${error.statusText}`;
      }
      return message;
    },

    /**
     * Triggers the hidden file input click.
     */
    changeAvatar() {
      if (this.isSaving) return;
      this.$refs.avatarInput?.click();
    },

    /**
     * Handles the file selection event for the avatar.
     * @param {Event} event - The file input change event.
     */
    handleAvatarChange(event) {
      const file = event.target.files?.[0];
      if (!file) return;

      console.log("UserProfileEdit: Avatar file selected:", file);
      // TODO: Implement avatar upload logic:
      // 1. Maybe show a preview using URL.createObjectURL(file).
      // 2. Store the file: this.selectedAvatarFile = file;
      // 3. Update saveProfile logic to handle file upload (likely FormData & separate API call).
      alert("Загрузка аватара пока не реализована.");

      // Reset input value to allow selecting the same file again
      if (this.$refs.avatarInput) {
        this.$refs.avatarInput.value = "";
      }
    },

    /**
     * Handles errors when loading the avatar image.
     * @param {Event} event - The error event from the img tag.
     */
    onAvatarError(event) {
      console.warn("UserProfileEdit: Failed to load avatar image.", event);
      // Optionally handle the error, e.g., hide the broken image
      // event.target.style.display = 'none'; // Hide broken image
    },

    /**
     * Focuses the input element associated with the given ID.
     * Uses $refs for reliable DOM access.
     * @param {string} inputId - The ID of the input/textarea to focus.
     */
    focusInput(inputId) {
      if (this.isSaving) return;
      const container = this.$refs.formContainer;
      if (container) {
        const inputElement = container.querySelector(`#${inputId}`);
        if (inputElement) {
          inputElement.focus();
        } else {
          console.warn(
            `UserProfileEdit: Element with ID #${inputId} not found within form container.`
          );
        }
      } else {
        console.warn(
          "UserProfileEdit: Form container ref not found. Ensure <fieldset ref='formContainer'> exists."
        );
      }
    },

    /**
     * Adjusts the height of the bio textarea to fit its content.
     */
    adjustTextareaHeight() {
      const textarea = this.$refs.bioTextarea;
      if (textarea) {
        textarea.style.height = "auto";
        const scrollHeight = textarea.scrollHeight;
        const newHeight = scrollHeight + 2;
        textarea.style.height = `${newHeight}px`;
      }
    },
  },
};
</script>

<style scoped>
@import "@/styles/userprofile.css";
@import "@/styles/userprofileedit.css";
</style>
