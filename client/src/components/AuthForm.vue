<template>
  <div class="auth-container">
    <!-- Логотип в левом нижнем углу -->
    <div class="logo-corner">&laquo;GoChat©&raquo; 2025</div>
    <!-- Кнопка смены темы в правом нижнем углу -->
    <button
      class="theme-switch-corner"
      @click="$emit('toggle-theme')"
      aria-label="Переключить тему"
    >
      <svg
        class="theme-icon"
        :class="{ 'animate-icon': animateIcon }"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <path :d="isDarkTheme ? moonIcon : sunIcon" />
      </svg>
    </button>

    <transition name="fade-slide">
      <div class="auth-card">
        <div class="auth-header">
          <!-- Центрированный заголовок -->
          <h2>{{ isLoginMode ? "Вход" : "Регистрация" }}</h2>
        </div>

        <!-- Форма авторизации -->
        <form @submit.prevent="handleSubmit" novalidate>
          <!-- Группа для Email -->
          <div
            class="auth-input-group"
            :style="{ marginBottom: errors.email ? '2.5rem' : '1.5rem' }"
          >
            <input
              v-model="email"
              type="email"
              placeholder=" "
              class="auth-input"
            />
            <label
              class="auth-input-label"
              :class="{ 'error-label': errors.email }"
            >
              Почта
            </label>
            <transition name="fade">
              <div v-if="errors.email" class="error-message">
                {{ errors.email }}
              </div>
            </transition>
          </div>

          <!-- Группа для Пароля -->
          <div
            class="auth-input-group"
            :style="{ marginBottom: errors.password ? '2.5rem' : '1.5rem' }"
          >
            <input
              v-model="password"
              :type="showPassword ? 'text' : 'password'"
              placeholder=" "
              class="auth-input"
            />
            <label
              class="auth-input-label"
              :class="{ 'error-label': errors.password }"
            >
              Пароль
            </label>
            <button
              type="button"
              class="password-toggle"
              @click="showPassword = !showPassword"
            >
              <svg
                class="eye-icon"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <path :d="showPassword ? eyeOpenIcon : eyeClosedIcon" />
              </svg>
            </button>
            <transition name="fade">
              <div v-if="errors.password" class="error-message">
                {{ errors.password }}
              </div>
            </transition>
          </div>

          <!-- Кнопка отправки -->
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
import { ref, defineProps, defineEmits, watch } from "vue";

const props = defineProps({
  isDarkTheme: { type: Boolean, default: false },
});
const emit = defineEmits(["login", "toggle-theme"]);

const isLoginMode = ref(true);
const email = ref("");
const password = ref("");
const showPassword = ref(false);

const errors = ref({
  email: "",
  password: "",
});

/**
 * Иконки в стиле heroicons
 */
const sunIcon = `
  M12 3v1
  m0 16v1
  m9-9h-1
  M4 12H3
  m15.364-6.364l-.707.707
  M6.343 17.657l-.707.707
  m12.728 0l-.707-.707
  M6.343 6.343l-.707-.707
  M12 8a4 4 0 110 8 4 4 0 010-8z
`;
const moonIcon = `
  M21 12.79A9 9 0 1111.21 3
  a7 7 0 009.79 9.79z
`;
const eyeOpenIcon = `
  M1 12c1.67-4.33 5.477-7 11-7
  s9.33 2.67 11 7
  c-1.67 4.33-5.477 7-11 7
  s-9.33-2.67-11-7z
  M12 15a3 3 0 100-6 3 3 0 000 6z
`;
const eyeClosedIcon = `
  M13.875 18.825A10.05 10.05 0 0112 19
  c-4.478 0-8.373-2.946-9.542-7
  a10.05 10.05 0 012.043-3.95
  m2.004-2.093A9.974 9.974 0 0112 5
  c4.478 0 8.373 2.946 9.542 7
  -.225.78-.52 1.529-.885 2.234
  M15 12a3 3 0 01-3 3
  m0-6a3 3 0 013 3
  M3 3l18 18
`;

const animateIcon = ref(false);
watch(
  () => props.isDarkTheme,
  () => {
    animateIcon.value = true;
    setTimeout(() => {
      animateIcon.value = false;
    }, 600);
  }
);

const handleSubmit = () => {
  errors.value.email = "";
  errors.value.password = "";

  if (!email.value.trim()) {
    errors.value.email = "Пожалуйста, введите Почту.";
  } else {
    const emailRegex = /.+@.+\..+/;
    if (!emailRegex.test(email.value)) {
      errors.value.email = "Неверный формат Почты.";
    }
  }
  if (!password.value.trim()) {
    errors.value.password = "Пожалуйста, введите пароль.";
  }
  if (!errors.value.email && !errors.value.password) {
    try {
      emit("login");
    } catch (error) {
      console.error("Auth error:", error);
    }
  }
};
</script>

<style scoped>
@import "@/styles/variables.css";
@import "@/styles/chat.css";

/* Общие стили для контейнера */
.auth-container {
  user-select: none;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background: var(--primary-light);
  transition: background 0.5s ease;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  overflow-y: auto;
}

/* Логотип в левом нижнем углу */
.logo-corner {
  position: absolute;
  bottom: 10px;
  left: 10px;
  font-size: 1.1rem;
  font-weight: bold;
  color: rgba(0, 0, 0, 0.5);
  z-index: 10;
}
.dark-theme .logo-corner {
  color: rgba(255, 255, 255, 0.3);
}

/* Кнопка смены темы в правом нижнем углу */
.theme-switch-corner {
  position: absolute;
  bottom: 10px;
  right: 10px;
  background: none;
  border: none;
  padding: 0.5rem;
  cursor: pointer;
  z-index: 10;
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.theme-switch-corner .theme-icon {
  width: 32px;
  height: 32px;
  transition: transform 0.6s ease, stroke 0.3s ease;
  stroke: rgba(0, 0, 0, 0.5);
}
.dark-theme .theme-switch-corner .theme-icon {
  stroke: rgba(255, 255, 255, 0.3);
}
.animate-icon {
  animation: icon-pop 0.6s ease;
}
@keyframes icon-pop {
  0% {
    transform: scale(1) rotate(0deg);
  }
  50% {
    transform: scale(1.2) rotate(45deg);
  }
  100% {
    transform: scale(1) rotate(0deg);
  }
}

.auth-card {
  width: 100%;
  max-width: 400px;
  padding: 2.5rem;
  background: var(--glass-effect-light, rgba(255, 255, 255, 0.9));
  backdrop-filter: blur(15px);
  border-radius: 1.5rem;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  transition: all 0.5s ease;
}

/* Сохранение тёмной темы для десктопа */
.dark-theme .auth-container {
  background: var(--primary-dark);
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
  text-align: center;
  margin-bottom: 2rem;
}
.auth-header h2 {
  font-size: 1.75rem;
  font-weight: 600;
  color: var(--text-light);
  transition: color 0.5s ease;
  margin: 0;
}
.dark-theme .auth-header h2 {
  color: var(--text-dark);
}

.auth-input-group {
  position: relative;
  transition: margin-bottom 0.3s ease;
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
.auth-input:focus {
  outline: none;
  box-shadow: inset 0 0 5px rgba(0, 0, 0, 0.1);
}
.dark-theme .auth-input {
  background: rgba(255, 255, 255, 0.05);
  color: var(--text-dark);
}

.auth-input-label {
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
.dark-theme .auth-input-label {
  color: var(--text-dark);
}
.auth-input:focus ~ .auth-input-label,
.auth-input:not(:placeholder-shown) ~ .auth-input-label {
  top: 0;
  font-size: 0.875rem;
  color: var(--gradient-sent-end);
}
.error-label {
  color: #ff4f4f;
}
.dark-theme .error-label {
  color: #ff4f4f;
}

.error-message {
  position: absolute;
  left: 1rem;
  bottom: -1.2rem;
  font-size: 0.875rem;
  color: #ff4f4f;
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(-5px);
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
  width: 24px;
  height: 24px;
  transition: stroke 0.5s ease;
}

.auth-button {
  width: 100%;
  padding: 1rem;
  background: var(--gradient-sent-end);
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
  color: var(--gradient-sent-end);
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
  transition: color 0.5s ease;
}
.dark-theme .forgot-password {
  color: var(--text-dark);
}

/* Мобильная адаптация */
@media (max-width: 480px) {
  .auth-container {
    align-items: center;
    padding: 1rem;
    height: 100vh;
  }
  /* Убираем контейнерную стилизацию для формы (сохраняем только содержимое) */
  .auth-card {
    width: 100%;
    max-width: 100%;
    padding: 1rem;
    margin-top: 1rem;
    background: transparent;
    box-shadow: none;
    border-radius: 0;
    transform: none;
  }
  .auth-header h2 {
    font-size: 1.5rem;
  }
  .auth-input {
    font-size: 0.9rem;
    padding: 0.8rem;
  }
  .auth-button {
    padding: 0.8rem;
    font-size: 0.9rem;
  }
  /* Фиксируем логотип и кнопку смены темы вне прокрутки */
  .logo-corner,
  .theme-switch-corner {
    position: fixed;
  }
  .logo-corner {
    bottom: 10px;
    left: 10px;
    font-size: 0.9rem;
  }
  .theme-switch-corner {
    bottom: 5px;
    right: 10px;
    width: 40px;
    height: 40px;
  }
  .theme-switch-corner .theme-icon {
    width: 28px;
    height: 28px;
  }
  /* Если используется тёмная тема, убираем контейнерное оформление формы */
  .dark-theme .auth-card {
    background: transparent;
    box-shadow: none;
    border-radius: 0;
  }
}

.dark-theme .theme-icon,
.dark-theme .eye-icon {
  color: white;
}
</style>
