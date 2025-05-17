<template>
    <div class="text-editor">
      
      <!-- Панель инструментов (видна только в режиме редактирования) -->
      <div v-if="isEditing && editor" class="toolbar">
        <button @click="editor.chain().focus().toggleBold().run()" :class="{ 'is-active': editor.isActive('bold') }">
          Bold
        </button>
        <button @click="editor.chain().focus().toggleItalic().run()" :class="{ 'is-active': editor.isActive('italic') }">
          Italic
        </button>
        <button @click="editor.chain().focus().toggleStrike().run()" :class="{ 'is-active': editor.isActive('strike') }">
          Strike
        </button>
        <button @click="editor.chain().focus().toggleBulletList().run()" :class="{ 'is-active': editor.isActive('bulletList') }">
          Bullet List
        </button>
        <button @click="editor.chain().focus().toggleOrderedList().run()" :class="{ 'is-active': editor.isActive('orderedList') }">
          Ordered List
        </button>
        <button @click="editor.chain().focus().toggleBlockquote().run()" :class="{ 'is-active': editor.isActive('blockquote') }">
          Blockquote
        </button>
        <button @click="editor.chain().focus().undo().run()">
          Undo
        </button>
        <button @click="editor.chain().focus().redo().run()">
          Redo
        </button>
      </div>
  
      <!-- Редактор Tiptap -->
      <div v-if="editor" class="editor">
        <editor-content :editor="editor" />
      </div>
  
        <!-- Кнопки "Опубликовать" и "Редактировать" -->
        <div class="buttons-block">
            <button @click="toggleEditMode" class="edit-button">
                {{ isEditing ? 'Скрыть панель' : 'Редактировать' }}
            </button>
        
            <button @click="publish" class="publish-button">Опубликовать</button>
        </div>
    </div>
  </template>
  
  <script>
  import { Editor, EditorContent } from '@tiptap/vue-3'
  import StarterKit from '@tiptap/starter-kit'
  
  export default {
    components: {
      EditorContent,
    },
    data() {
      return {
        editor: null, // Экземпляр редактора
        isEditing: false, // Режим редактирования 
      }
    },
    mounted() {
      // Инициализация редактора
      this.editor = new Editor({
        content: '<p>Напишите что-нибудь...</p>',
        extensions: [StarterKit],
      })
    },
    beforeUnmount() {
      // Уничтожение редактора
      this.editor.destroy()
    },
    methods: {
      // Переключение режима редактирования
      toggleEditMode() {
        this.isEditing = !this.isEditing
      },
      // Логика для кнопки "Опубликовать"
      publish() {
        const content = this.editor.getHTML()
        console.log('Опубликовано:', content)
        alert('Сообщение опубликовано! Проверьте консоль.')
      },
    },
  }
  </script>
  
  <style scoped>
  .text-editor {
    max-width: 800px;
    margin: 20px auto;
    padding: 20px;
    border: 1px solid #ccc;
    border-radius: 8px;
    background-color: #f9f9f9;
  }
  
  .toolbar {
    margin-bottom: 10px;
  }
  
  .toolbar button {
    background-color: #f1f1f1;
    border: none;
    padding: 8px 12px;
    margin: 4px;
    border-radius: 4px;
    cursor: pointer;
  }
  
  .toolbar button.is-active {
    background-color: #42b983;
    color: white;
  }
  
  .editor {
    margin-bottom: 20px;
  }

  :deep(.ProseMirror) {
    background-color: white;
    padding: 10px;
  }
  
  .buttons-block {
    display: flex;

    justify-content: space-between;
  }

  .edit-button,
  .publish-button {
    padding: 10px 20px;
    background-color: #42b983;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 16px;
    margin-right: 10px;
  }
  
  .edit-button:hover,
  .publish-button:hover {
    background-color: #3aa876;
  }


  </style>

/*
Установила и подключила TipTap, тк она без проблем интегрируется в проекты Vue 
и является высокоуровневой (в отличие от ProseMirror, на которой основана)

Почитала документацию (про основные команды и цепочки команд)
Узнала, как обратиться в CSS к содержимому экземпляра редактора
Узнала, что из-за динамического добавления 
класса ProseMirror внутри ContentEditor нельзя 
поменять стили элементов этого класса или дочерних внутри компонента, 
использующего ContentEditor. Поэтому такие стили нужно либо выносить 
в global.css, либо использовать глубокий селектор
*/