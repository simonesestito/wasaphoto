import {createRouter, createWebHashHistory} from 'vue-router';
import HomeView from '../views/HomeView.vue';
import MyAccountView from '../views/MyAccountView.vue';
import SearchUsers from '../views/SearchUsers.vue';

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: HomeView},
		{path: '/me', component: MyAccountView},
		{path: '/search', component: SearchUsers},
	]
})

export default router
