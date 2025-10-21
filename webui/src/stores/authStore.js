// stores/authStore.js
import { reactive } from "vue";
import { doLogin } from "../services/auth";
import { LocalStorage } from "../config/localStorage";

const auth = reactive({
	userId: localStorage.getItem(LocalStorage.userId) || null,
});

async function login(username) {
	const res = await doLogin(username);

	if (res.status === 200 || res.status === 201) {
		auth.userId = res.data.id;
		localStorage.setItem(LocalStorage.userId, res.data.id);

		return { ok: true };
	}

	return {
		ok: false,
		error: res?.error || res.response?.data?.message || res.response?.data,
	};
}

export { auth, login };
