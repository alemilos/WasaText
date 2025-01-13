import { createRouter, createWebHashHistory } from "vue-router";

// Pages
import HomeView from "../views/HomeView.vue";
import LoginView from "../views/LoginView.vue";
import { clearUserId } from "../services/auth";

const router = createRouter({
	history: createWebHashHistory("/"),
	routes: [
		{ path: "/", component: HomeView },
		{ path: "/login", component: LoginView },
	],
});

/**
 * Redirect the user to login force him to logout  
 */
router.logoutAndRedirect = async function () {
	clearUserId() // remove user id 

	await this.replace({
		path: '/login',
		query: { previous: router.currentRoute.value.path },
	});
}.bind(router) 

export default router;
