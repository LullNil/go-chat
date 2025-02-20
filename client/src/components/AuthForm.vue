<template>
  <div class="auth-container">
    <transition name="fade-slide">
      <div class="auth-card" :class="{ 'dark-theme': isDark }">
        <div class="auth-header">
          <h2>{{ isLoginMode ? "Вход" : "Регистрация" }}</h2>
          <button
            class="theme-toggle"
            @click="toggleTheme"
            aria-label="Переключить тему"
          >
            <svg class="theme-icon" viewBox="0 0 24 24">
              <path :d="isDark ? moonIcon : sunIcon" />
            </svg>
          </button>
        </div>

        <form @submit.prevent="handleSubmit">
          <div class="input-group">
            <input
              v-model="email"
              type="email"
              required
              placeholder=" "
              class="auth-input"
            />
            <label class="input-label">Email</label>
          </div>

          <div class="input-group">
            <input
              v-model="password"
              :type="showPassword ? 'text' : 'password'"
              required
              placeholder=" "
              class="auth-input"
            />
            <label class="input-label">Пароль</label>
            <button
              type="button"
              class="password-toggle"
              @click="showPassword = !showPassword"
            >
              <svg class="eye-icon" viewBox="0 0 24 24">
                <path :d="showPassword ? eyeOpenIcon : eyeClosedIcon" />
              </svg>
            </button>
          </div>

          <button type="submit" class="auth-button">
            {{ isLoginMode ? "Войти" : "Зарегистрироваться" }}
          </button>
        </form>

        <div class="auth-footer">
          <button class="mode-toggle" @click="isLoginMode = !isLoginMode">
            {{
              isLoginMode ? "Нет аккаунта? Создать" : "Уже есть аккаунт? Войти"
            }}
          </button>
          <a href="#" class="forgot-password">Забыли пароль?</a>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, defineEmits } from "vue";

const emit = defineEmits(["login"]); // Объявите событие

const isDark = ref(false);
const isLoginMode = ref(true);
const email = ref("");
const password = ref("");
const showPassword = ref(false);

const sunIcon =
  "M12 18a6 6 0 1 0 0-12 6 6 0 0 0 0 12zM11 1h2v3h-2V1zm0 19h2v3h-2v-3zM3.515 4.929l1.414-1.414L7.05 5.636 5.636 7.05 3.515 4.93zM16.95 18.364l1.414-1.414 2.121 2.121-1.414 1.414-2.121-2.121zm2.121-14.85l1.414 1.415-2.121 2.121-1.414-1.414 2.121-2.121zM5.636 16.95l1.414 1.414-2.121 2.121-1.414-1.414 2.121-2.121z";
const moonIcon =
  "M12 3c.132 0 .263 0 .393 0a7.5 7.5 0 0 0 7.92 12.446a9 9 0 1 1-8.313-12.454z";
const eyeOpenIcon =
  "M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z M15 12a3 3 0 1 1-6 0 3 3 0 0 1 6 0z";
const eyeClosedIcon =
  "M11.83 9L15 12.16V12a3 3 0 0 0-3-3h-.17zm-4.3.8l1.55 1.55a3 3 0 0 0 4.12 4.12l1.55 1.55A9.96 9.96 0 0 1 3 12c1.65-3.66 5.22-6 9-6 1.65 0 3.19.42 4.57 1.14l-1.43 1.43A7.88 7.88 0 0 0 12 5c-2.97 0-5.58 1.5-7.47 3.8zM2.81 2.81L1.39 4.22 4.62 7.45C3.16 8.73 2 10.28 2 12c1.73 4.39 6 7.5 11 7.5 1.55 0 3.03-.3 4.38-.84l2.77 2.77 1.41-1.41L2.81 2.81z";

const toggleTheme = () => {
  isDark.value = !isDark.value;
  document.documentElement.classList.toggle("dark-theme", isDark.value);
};

const handleSubmit = async () => {
  try {
    emit("login"); // Только эмит события
  } catch (error) {
    console.error("Auth error:", error);
  }
};
</script>

<style scoped>
@import "@/styles/variables.css";
@import "@/styles/chat.css";

.auth-container {
  position: fixed; /* Добавлено */
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  display: flex; /* Добавлено */
  justify-content: center;
  align-items: center;
  background: var(--primary-light);
  z-index: 1000; /* Добавлено */
}

.dark-theme .auth-container {
  background: var(--primary-dark);
}

.auth-card {
  width: 100%;
  max-width: 400px;
  padding: 2.5rem;

  background: var(--glass-effect-light, rgba(255, 255, 255, 0.9));
  backdrop-filter: blur(15px);
  border-radius: 1.5rem;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.dark-theme .auth-card {
  background: var(--secondary-dark);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  border-color: rgba(255, 255, 255, 0.05);
}

.fade-slide-enter-active {
  transition: all 0.6s cubic-bezier(0.23, 1, 0.32, 1);
}

.fade-slide-enter-to {
  opacity: 1;
  transform: translateY(0);
}

.auth-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.auth-header h2 {
  font-size: 1.75rem;
  font-weight: 600;
  color: var(--text-light);
}

.dark-theme .auth-header h2 {
  color: var(--text-dark);
}

.theme-toggle {
  background: none;
  border: none;
  padding: 0.5rem;
  cursor: pointer;
  transition: transform 0.3s ease;
}

.theme-toggle:hover {
  transform: rotate(15deg);
}

.theme-icon {
  width: 24px;
  height: 24px;
  fill: currentColor;
}

.input-group {
  position: relative;
  margin-bottom: 1.5rem;
}

.auth-input {
  width: 100%;
  padding: 1rem;
  background: rgba(255, 255, 255, 0.1);
  border: none;
  border-radius: 0.75rem;
  font-size: 1rem;
  transition: all 0.3s ease;
  outline: none;
  box-shadow: inset 0 0 5px rgba(0, 0, 0, 0.1);
}

.dark-theme .auth-input {
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(255, 255, 255, 0.1);
  color: var(--text-dark);
}

.input-label {
  position: absolute;
  left: 1rem;
  top: 50%;
  transform: translateY(-50%);
  pointer-events: none;
  transition: all 0.3s ease;
  color: var(--text-light);
  background: none;
  padding: 0 0.5rem;
}

.dark-theme .input-label {
  color: var(--text-dark);
  background: none;
}

.auth-input:focus ~ .input-label,
.auth-input:not(:placeholder-shown) ~ .input-label {
  top: 0;
  font-size: 0.875rem;
  color: var(--accent-blue);
}

.password-toggle {
  position: absolute;
  right: 1rem;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  cursor: pointer;
}

.eye-icon {
  width: 20px;
  height: 20px;
  fill: none;
  stroke: currentColor;
  stroke-width: 1.5;
}

.auth-button {
  width: 100%;
  padding: 1rem;
  background: var(--accent-blue);
  color: white;
  border: none;
  border-radius: 0.75rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.auth-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(88, 101, 242, 0.25);
}

.auth-footer {
  margin-top: 1.5rem;
  text-align: center;
}

.mode-toggle {
  background: none;
  border: none;
  color: var(--accent-blue);
  cursor: pointer;
  font-weight: 500;
  transition: opacity 0.3s ease;
}

.mode-toggle:hover {
  opacity: 0.8;
}

.forgot-password {
  display: block;
  margin-top: 1rem;
  color: var(--text-light);
  text-decoration: none;
  font-size: 0.875rem;
}

.dark-theme .forgot-password {
  color: var(--text-dark);
}

@media (max-width: 480px) {
  .auth-card {
    padding: 1.5rem;
    margin: 1rem;
  }
}
</style>
