import {createRouter, createWebHashHistory} from 'vue-router';
import HomeView from '../views/HomeView.vue';
import MyAccountView from '../views/MyAccountView.vue';
import SearchUsers from '../views/SearchUsers.vue';
import LoginView from "../views/LoginView.vue";
import {getCurrentUID} from "../services/auth-store";
import EditAccountView from "../views/EditAccountView.vue";
import SingleUserView from "../views/SingleUserView.vue";
import FollowersView from "../views/FollowersView.vue";
import FollowingsView from "../views/FollowingsView.vue";

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: HomeView},
		{path: '/me', component: MyAccountView},
		{path: '/search', component: SearchUsers},
		{path: '/login', component: LoginView},
		{path: '/me/edit', component: EditAccountView},
		{path: '/users/:username', component: SingleUserView},
		{path: '/users/:username/followers', component: FollowersView},
		{path: '/users/:username/followings', component: FollowingsView},
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
