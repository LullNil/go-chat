<template>
  <!-- Auth form -->
  <form @submit.prevent="handleSubmit" novalidate>
    <!-- Username field -->
    <div class="extra-field username-field" :class="{ active: !isLoginMode }">
      <div class="auth-input-group with-icon">
        <input
          v-model="username"
          type="text"
          placeholder=" "
          class="auth-input"
          @input="clearError('username')"
        />
        <label
          class="auth-input-label"
          :class="{ 'error-label': errors.username }"
        >
          Имя пользователя
        </label>

        <!-- Username tooltip -->
        <UsernameTooltip />

        <!-- Error messages -->
        <ErrorMessage :message="errors.username" />
      </div>
    </div>

    <!-- Email or username field -->
    <div class="auth-input-group">
      <input
        v-model="email"
        type="text"
        placeholder=" "
        class="auth-input"
        @input="clearError('email')"
      />
      <label class="auth-input-label" :class="{ 'error-label': errors.email }">
        {{ isLoginMode ? "Почта или логин" : "Почта" }}
      </label>

      <!-- Error messages -->
      <ErrorMessage :message="generalServerError" />
      <ErrorMessage :message="errors.email" />
    </div>

    <!-- Password field -->
    <div class="auth-input-group with-icon">
      <input
        v-model="password"
        :type="showPassword ? 'text' : 'password'"
        placeholder=" "
        class="auth-input"
        @input="clearError('password')"
      />
      <label
        class="auth-input-label"
        :class="{ 'error-label': errors.password }"
      >
        Пароль
      </label>

      <!-- Password toggle button -->
      <button type="button" class="password-toggle" @click="togglePassword">
        <svg
          class="eye-icon"
          :class="{ 'animate-closed': !showPassword }"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="1.5"
          stroke-linecap="butt"
          stroke-linejoin="round"
        >
          <path
            d="M1 12c1.67-4.33 5.477-7 11-7s9.33 2.67 11 7c-1.67 4.33-5.477 7-11 7s-9.33-2.67-11-7zM12 15a3 3 0 100-6 3 3 0 000 6z"
          />
          <line
            class="eye-cross"
            :class="{ drawn: !showPassword }"
            x1="4"
            y1="4"
            x2="20"
            y2="20"
          />
        </svg>
      </button>

      <!-- Error messages -->
      <ErrorMessage :message="errors.password" />
    </div>

    <!-- Confirm password field (for registration mode only) -->
    <div class="extra-field confirm-field" :class="{ active: !isLoginMode }">
      <div class="auth-input-group">
        <input
          v-model="confirmPassword"
          type="password"
          placeholder=" "
          class="auth-input"
          @input="clearError('confirmPassword')"
        />
        <label
          class="auth-input-label"
          :class="{ 'error-label': errors.confirmPassword }"
        >
          Подтвердите пароль
        </label>

        <!-- Error messages -->
        <ErrorMessage :message="errors.confirmPassword" />
      </div>
    </div>

    <!-- Auth button -->
    <button type="submit" class="auth-button">
      {{ isLoginMode ? "Войти" : "Зарегистрироваться" }}
    </button>
  </form>
</template>

<script>
import ErrorMessage from "./ErrorMessage.vue";
import UsernameTooltip from "./UsernameTooltip.vue";

export default {
  name: "AuthForm",
  components: {
    ErrorMessage,
    UsernameTooltip,
  },
  props: {
    isLoginMode: {
      type: Boolean,
      required: true,
    },
    serverError: {
      type: String,
      default: "",
    },
  },
  emits: ["auth-submit", "update:errorMessage", "update:serverError"],
  data() {
    return {
      username: "",
      email: "",
      password: "",
      confirmPassword: "",
      showPassword: false,
      errors: {
        username: "",
        email: "",
        password: "",
        confirmPassword: "",
      },
      errorMap: {
        "user already exists": "Такой пользователь уже существует.",
        "invalid email or password": "Неверный логин или пароль.",
        // Additionally add other error messages here
      },
    };
  },
  computed: {
    /**
     * Returns the translated server error message if it's
     * present in this.errorMap,
     * otherwise returns the original server error message.
     *
     * @returns {string} translated error message or original error message
     */
    generalServerError() {
      if (this.serverError) {
        const lowerError = this.serverError.toLowerCase();
        for (const key in this.errorMap) {
          if (lowerError.includes(key)) {
            return this.errorMap[key];
          }
        }
        return this.serverError;
      }
      return "";
    },
  },
  watch: {
    isLoginMode(newMode, oldMode) {
      if (newMode !== oldMode) {
        this.resetForm();
      }
    },
  },
  methods: {
    // resetErrors clears all error messages
    resetErrors() {
      this.errors = {
        username: "",
        email: "",
        password: "",
        confirmPassword: "",
      };
    },
    // clearError emits an event to clear the server error message
    clearError(field) {
      this.errors[field] = "";
      if (this.serverError) {
        this.$emit("update:errorMessage", "");
        this.$emit("update:serverError", "");
      }
    },
    // resetForm resets all form fields
    resetForm() {
      this.username = "";
      this.email = "";
      this.password = "";
      this.confirmPassword = "";
      this.resetErrors();
    },
    /**
     * Validates the form data based on the current mode (login or register)
     *
     * @returns {boolean} whether the form data is valid or not
     */
    validateForm() {
      let hasErrors = false;
      this.resetErrors();

      if (this.isLoginMode) {
        const loginInput = this.email.trim();
        if (!loginInput) {
          this.errors.email = "Введите почту или логин.";
          hasErrors = true;
        } else if (loginInput.includes("@")) {
          const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
          if (!emailRegex.test(loginInput)) {
            this.errors.email = "Неверный формат почты.";
            hasErrors = true;
          }
        }
        if (!this.password.trim()) {
          this.errors.password = "Введите пароль.";
          hasErrors = true;
        }
      } else {
        const emailValue = this.email.trim();
        if (!emailValue) {
          this.errors.email = "Введите почту.";
          hasErrors = true;
        } else {
          const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
          if (!emailRegex.test(emailValue)) {
            this.errors.email = "Неверный формат почты.";
            hasErrors = true;
          }
        }

        this.username = this.username.toLowerCase().trim();
        if (!this.username) {
          this.errors.username = "Введите имя пользователя.";
          hasErrors = true;
        } else if (!/^[a-z0-9_]{5,}$/.test(this.username)) {
          this.errors.username = "Недопустимое имя пользователя.";
          hasErrors = true;
        } else if (this.username.length > 30) {
          this.errors.username = "Слишком длинное имя пользователя.";
          hasErrors = true;
        }

        if (!this.password.trim()) {
          this.errors.password = "Введите пароль.";
          hasErrors = true;
        }

        if (this.password !== this.confirmPassword) {
          this.errors.confirmPassword = "Пароли не совпадают.";
          hasErrors = true;
        }
      }
      return !hasErrors;
    },
    /**
     * Handles a server error by translating it if it's present in
     * this.errorMap
     * and emits an event to update the error message.
     *
     * @param {string} serverError the server error message
     */
    handleServerError(serverError) {
      const lowerError = serverError.toLowerCase();
      for (const key in this.errorMap) {
        if (lowerError.includes(key)) {
          this.$emit("update:errorMessage", this.errorMap[key]);
          return;
        }
      }
      this.$emit("update:errorMessage", serverError);
    },
    /**
     * Handles form submission by validating the input fields and preparing
     * the authentication data.
     *
     * If the form is invalid, emits an event to update the error message
     * and signals submission failure.
     * If the form is valid, constructs the authData object based on
     * the current mode (login or register),
     * and emits an event with the authentication data.
     */
    handleSubmit() {
      if (!this.validateForm()) {
        const firstError = Object.values(this.errors).find((error) => error);
        this.$emit("update:errorMessage", firstError);
        this.$emit("auth-submit", false);
        return;
      }

      let authData = {
        password: this.password,
      };

      if (this.isLoginMode) {
        const loginInput = this.email.trim();
        if (loginInput.includes("@")) {
          authData.email = loginInput;
          authData.username = "";
        } else {
          authData.email = "";
          authData.username = loginInput;
        }
      } else {
        authData.email = this.email.trim();
        authData.username = this.username.trim();
        authData.confirmPassword = this.confirmPassword;
      }
      this.$emit("auth-submit", authData);
    },
    // togglePassword toggles the visibility of the password
    togglePassword() {
      this.showPassword = !this.showPassword;
    },
  },
};
</script>

<style scoped>
@import "@/styles/authpage.css";
</style>
