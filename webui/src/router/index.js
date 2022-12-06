import {createRouter, createWebHashHistory} from 'vue-router';
import HomeView from '../views/HomeView.vue';
import MyAccountView from '../views/MyAccountView.vue';
import SearchUsers from '../views/SearchUsers.vue';
import LoginView from "../views/LoginView.vue";
import {getCurrentUID} from "../services/auth-store";

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: HomeView},
		{path: '/me', component: MyAccountView},
		{path: '/search', component: SearchUsers},
		{path: '/login', component: LoginView},
	]
});

router.beforeEach((to, from, next) => {
	const goingToAnonymous = to.path === '/login';
	const isAuthenticated = getCurrentUID();
	if (goingToAnonymous || isAuthenticated) {
		// Allow routing
		next();
	} else {
		// Redirect to log in page
		next('/login');
	}
});

export default router
