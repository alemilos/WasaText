import { createRouter, createWebHashHistory } from "vue-router";

// Pages
import HomeView from "../views/HomeView.vue";
import LoginView from "../views/LoginView.vue";
import { clearUserId, getUserId } from "../services/auth";

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
	clearUserId(); // remove user id

	await this.replace({
		path: "/login",
		query: { previous: router.currentRoute.value.path },
	});
}.bind(router);

/**
 * Make sure the user is logged in (something like PrivateRouter)
 */
router.beforeEach((to, from, next) => {
	const userId = getUserId();

	if (!userId && to.path !== "/login") {
		// if no user and trying to access protected page
		next("/login");
	} else if (userId && to.path === "/login") {
		// if already logged in, block access to login
		next("/");
	} else {
		next(); // allow navigation
	}
});

export default router;
