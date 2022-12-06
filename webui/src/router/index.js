import {createRouter, createWebHashHistory} from 'vue-router';
import HomeView from '../views/HomeView.vue';
import MyAccountView from '../views/MyAccountView.vue';
import SearchUsers from '../views/SearchUsers.vue';
import LoginView from "../views/LoginView.vue";

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: HomeView},
		{path: '/me', component: MyAccountView},
		{path: '/search', component: SearchUsers},
		{path: '/login', component: LoginView},
	]
})

export default router
