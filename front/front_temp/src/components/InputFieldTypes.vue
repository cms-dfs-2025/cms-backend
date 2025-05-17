<template>
  <div class="input-field" :class="fieldModifiers">
    <!-- Стандартные поля и textarea -->
    <template v-if="isStandardInput || type === 'textarea'">
      <component
        :is="type === 'textarea' ? 'textarea' : 'input'"
        class="input-field__native"
        :type="type"
        :id="id"
        :value="modelValue"
        @input="handleInput"
        @focus="handleFocus"
        @blur="handleBlur"
        :placeholder="floatingLabel ? ' ' : placeholder"
        :disabled="disabled"
        :rows="type === 'textarea' ? rows : null"
        :class="{'input-error': localErrorMessage}"
      />
      <label 
        v-if="floatingLabel" 
        :for="id" 
        class="input-field__label" 
        :class="{
          'label-error': localErrorMessage,
          'label-hidden': !withAnimation && modelValue
        }"
      >
        {{ label }}
      </label>
    </template>
    
    <!-- File Input -->
    <template v-else-if="type === 'file'">
    <div class="input-field__file-wrapper">
      <input
        class="input-field__native input-field__native--file"
        type="file"
        :id="id"
        @change="handleFileChange"
        @focus="handleFocus"
        @blur="handleBlur"
        :disabled="disabled"
        :accept="accept"
      />
      <label :for="id" class="input-field__file-content" :class="{
        'input-error': localErrorMessage,
        'has-file': modelValue && !localErrorMessage
      }">
        <span v-if="localErrorMessage" class="file-error-message">{{ localErrorMessage }}</span>
        <span v-else class="input-field__file-text">
          {{ modelValue ? modelValue.name : fileLabel }}
        </span>
        <span class="input-field__file-icon" :class="{'icon-error': localErrorMessage}">
          <svg viewBox="0 0 24 24" width="24" height="24">
            <path v-if="modelValue" d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"/>
            <path v-else d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
          </svg>
        </span>
      </label>
      <div class="input-field__file-hint">1 файл, формата doc, docx, pdf</div>
    </div>
  </template>

    <!-- Dropdown -->
    <template v-else-if="type === 'dropdown'">
      <div class="input-field__dropdown-wrapper">
        <select
          class="input-field__native input-field__native--dropdown"
          :id="id"
          :value="modelValue"
          @change="handleInput"
          @focus="handleFocus"
          @blur="handleBlur"
          @click="toggleDropdown"
          :disabled="disabled"
          :class="{'input-error': localErrorMessage}"
        >
          <option v-if="placeholder" value="" disabled>
            {{ placeholder }}
          </option>
          <option
            v-for="option in options"
            :key="option.value"
            :value="option.value"
          >
            {{ option.label }}
          </option>
        </select>
        <span class="input-field__dropdown-icon" :class="{'is-open': isDropdownOpen}">
          <svg viewBox="0 0 24 24" width="24" height="24">
            <path v-if="isDropdownOpen" d="M12 8l-6 6 1.41 1.41L12 10.83l4.59 4.58L18 14z"/>
            <path v-else d="M16.59 8.59L12 13.17 7.41 8.59 6 10l6 6 6-6z"/>
          </svg>
        </span>
      </div>
    </template>

    <!-- Иконка справа для стандартных полей -->
    <span
      v-if="iconRight && (isStandardInput || type === 'textarea')"
      class="input-field__icon-right"
      @click="$emit('icon-right-click')"
    >
      <slot name="icon-right">
        <component :is="iconRight" />
      </slot>
    </span>

    <!-- Сообщение об ошибке (кроме file) -->
    <div v-if="localErrorMessage && type !== 'file'" class="error-message">{{ localErrorMessage }}</div>
  </div>
</template>

<script>
export default {
  name: 'InputFieldTypes',
  props: {
    type: {
      type: String,
      default: 'text',
      validator: function(value) {
        return [
          'text',
          'email',
          'password',
          'textarea',
          'file',
          'dropdown'
        ].includes(value);
      }
    },
    modelValue: {
      type: [String, Number, File, Object],
      default: ''
    },
    id: {
      type: String,
      required: true
    },
    label: {
      type: String,
      default: ''
    },
    placeholder: {
      type: String,
      default: ''
    },
    disabled: {
      type: Boolean,
      default: false
    },
    floatingLabel: {
      type: Boolean,
      default: true
    },
    withAnimation: {
      type: Boolean,
      default: true
    },
    rows: {
      type: Number,
      default: 3
    },
    options: {
      type: Array,
      default: function() {
        return [];
      }
    },
    fileLabel: {
      type: String,
      default: 'Прикрепить файл'
    },
    iconRight: {
      type: [String, Object],
      default: null
    },
    errorMessage: {
      type: String,
      default: ''
    },
    accept: {
      type: String,
      default: '.doc,.docx,.pdf'
    }
  },
  data() {
    return {
      isFocused: false,
      isDropdownOpen: false,
      localErrorMessage: ''
    };
  },
  watch: {
    errorMessage(newVal) {
      this.localErrorMessage = newVal;
    },
    modelValue() {
      // Сбрасываем ошибку при изменении значения
      this.localErrorMessage = '';
    }
  },
  computed: {
    isStandardInput: function() {
      return ['text', 'email', 'password'].includes(this.type);
    },
    fieldModifiers: function() {
      return {
        'input-field--disabled': this.disabled,
        'input-field--with-icon': !!this.iconRight,
        [`input-field--${this.type}`]: true,
        'with-animation': this.withAnimation
      };
    }
  },
  methods: {
    handleInput: function(event) {
      this.$emit('update:modelValue', event.target.value);
      this.localErrorMessage = ''; // Сбрасываем ошибку при вводе
    },
    handleFileChange: function(event) {
      const file = event.target.files[0];
      if (file) {
        // Проверка размера файла (15MB)
        if (file.size > 15 * 1024 * 1024) {
          this.localErrorMessage = 'Файл слишком большой (макс. 15MB)';
          this.$emit('update:modelValue', null);
          event.target.value = ''; // Сбрасываем значение input file
          //this.fileLabel = 'Файл слишком большой';
          return;
        }
        this.localErrorMessage = 'Прикрепить файл';
      }
      this.$emit('update:modelValue', file);
    },
    handleFocus: function(event) {
      this.isFocused = true;
      this.$emit('focus', event);
    },
    handleBlur: function(event) {
      this.isFocused = false;
      this.$emit('blur', event);
    },
    toggleDropdown: function() {
      this.isDropdownOpen = !this.isDropdownOpen;
    }
  },
  emits: [
    'update:modelValue',
    'focus',
    'blur',
    'icon-right-click'
  ]
};
</script>

<style scoped>
/* Подключение глобальных переменных */
.input-field {
  --input-text-color: #222;
  --input-label-color: rgba(34, 34, 34, 0.5);
  --input-border-hover: rgba(34, 34, 34, 0.5);
  --input-border-focus: #222;
  --input-error-color: #cb3d35;
  --input-warning-color: #cb3d35;
  font-family: var(--font-family-base);
  position: relative;
  margin-bottom: 1.5rem;
}

/* Базовые стили для всех input/textarea */
.input-field__native {
  width: 100%;
  height: 55px;
  padding: 18px 0 16px;
  box-sizing: border-box;
  font-family: inherit;
  font-weight: 350;
  font-size: 16px;
  color: var(--input-text-color);
  border: none;
  border-bottom: 1px solid var(--input-label-color);
  background: transparent;
  transition: border-color 0.3s ease;
}

.input-field__native:focus {
  outline: none;
  border-bottom-color: var(--input-border-focus);
  opacity: 1;
}

/* Стили для hover */
.input-field__native:not(.input-error):hover:not(:focus) {
  border-bottom-color: var(--input-border-hover);
}

/* Плавающий лейбл */
.input-field__label {
  position: absolute;
  top: 50%;
  left: 0;
  transform: translateY(-50%);
  pointer-events: none;
  font-size: 16px;
  color: var(--input-label-color);
  transition: all 150ms cubic-bezier(0.4, 0, 0.2, 1);
}

.input-field__native:hover + .input-field__label {
  opacity: 1;
}

.with-animation input:focus + .input-field__label,
.with-animation input:not(:placeholder-shown) + .input-field__label,
.with-animation textarea:focus + .input-field__label,
.with-animation textarea:not(:placeholder-shown) + .input-field__label {
  top: 0;
  transform: translateY(0);
  font-size: 12px;
  opacity: 0.5;
}

.label-simple {
  transition: none;
}

input:focus + .label-simple,
.label-simple.label-hidden {
  display: none;
}

.label-error {
  color: var(--input-error-color);
}

/* Специфичные стили для textarea */
.input-field__native--textarea {
  height: auto;
  min-height: 120px;
  padding-top: 18px;
  resize: vertical;
}

/* File input стили */
.input-field__file-wrapper {
  position: relative;
}

.input-field__native--file {
  opacity: 0;
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  border: 0;
}

.input-field__file-content {
  display: flex;
  justify-content: space-between;
  padding: 18px 0 16px;
  border-bottom: 1px solid var(--input-label-color);
  cursor: pointer;
  transition: border-color 0.3s ease;
  color: var(--input-label-color);
}

.input-field__file-text {
  flex-grow: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.file-error-message {
  color: var(--input-error-color);
  margin-left: 8px;
}

.input-field__file-content:hover {
  color: var(--input-text-color);
  border-bottom-color: var(--input-border-hover);
}

.input-field__file-content:hover .input-field__file-icon {
  fill: var(--input-text-color);
  opacity: 1;
}

.input-field__file-icon {
  fill: var(--input-label-color);
  transition: all 0.3s ease;
  margin-left: 8px;
}

.input-field__file-hint {
  font-size: 12px;
  color: var(--input-label-color);
  margin-top: 4px;
}

/* Dropdown стили */
.input-field__dropdown-wrapper {
  position: relative;
}

.input-field__native--dropdown {
  appearance: none;
  padding-right: 30px;
  color: var(--input-label-color);
  cursor: pointer;
}

.input-field__native--dropdown:hover {
  color: var(--input-text-color);
  border-bottom-color: var(--input-border-hover);
}

.input-field__native--dropdown:focus {
  color: var(--input-text-color);
  border-bottom-color: var(--input-border-focus);
  opacity: 1;
}

.input-field__dropdown-icon {
  position: absolute;
  right: 0;
  top: 50%;
  transform: translateY(-50%);
  fill: var(--input-label-color);
  pointer-events: none;
  transition: all 0.3s ease;
}

.input-field__native--dropdown:hover + .input-field__dropdown-icon {
  fill: var(--input-text-color);
  opacity: 1;
}

.input-field__dropdown-icon.is-open {
  fill: var(--input-text-color);
  opacity: 1;
}

/* Ошибки */
.input-error {
  color: var(--input-error-color);
  border-bottom-color: var(--input-error-color) !important;
}

.error-message {
  margin: 0;
  padding: 0;
  text-align: left;
  font-size: 12px;
  font-weight: 350;
  color: var(--input-error-color);
}

/* Иконка справа */
.input-field__icon-right {
  position: absolute;
  right: 0;
  top: 50%;
  transform: translateY(-50%);
  cursor: pointer;
  fill: var(--input-label-color);
  transition: color 0.3s ease;
}

.input-field__icon-right:hover {
  fill: var(--input-text-color);
  opacity: 1;
}

.input-field__file-content.has-file {
  color: var(--input-text-color);
  opacity: 1;
}

.input-field__file-icon.icon-error {
  fill: var(--input-error-color) !important;
}
</style>