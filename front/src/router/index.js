import { createRouter, createWebHistory } from 'vue-router';
import LoginForm from '../views/LoginForm.vue';
import TextEditor from '../views/TextEditor.vue';
/*import { compile } from 'vue';*/
import LoginForm2 from '@/views/LoginForm2.vue';
import ResetPassword from '@/views/ResetPassword.vue';
import NewPassword from '@/views/NewPassword.vue';
import RegForm from '@/views/RegForm.vue';

const routes = [
  {
    path: '/signup',
    name: 'SignUp',
    component: RegForm,
  },
  {
    path: '/test',
    name: 'Test',
    component: TextEditor,
  },
  {
    path: '/auth/login',
    name: 'Login',
    component: LoginForm2,
  },
  {
    path: '/auth/reset-password',
    name: 'Reset',
    component: ResetPassword,
  },
  {
    path: '/auth/create-password',
    name: 'Create',
    component: NewPassword,
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;