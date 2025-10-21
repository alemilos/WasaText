<script setup>
import "./style/LoginView.css";
import Logo from "../assets/images/logo.png";
import { hackerText } from "../utils/text";
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import { login } from "../stores/authStore";
import { useToast } from "vue-toast-notification";

const title = ref(null);
const username = ref("");
const router = useRouter();
const $toast = useToast();

async function submit(e) {
	e.preventDefault();

	const res = await login(username.value);
	if (res.ok) {
		router.push("/");
	} else {
		const error = "Impossibile effettuare il login";
		$toast.error(error);
	}
}

onMounted(() => {
	hackerText(title.value);
});
</script>

<template>
	<div class="login-container">
		<div class="form-container">
			<div class="form-heading">
				<img :src="Logo" alt="logo" />
				<h1 ref="title">Accedi a WasaText</h1>
				<p>
					Chatta in modo <span>sicuro</span> e <span>veloce</span>,
					con amici e familiari.
				</p>
			</div>

			<form @submit="submit">
				<input type="text" placeholder="username" v-model="username" />
				<button type="submit">Accedi</button>
			</form>

			<div class="form-info">
				<v-icon name="bi-info-circle" />
				<p>
					Fai attenzione, chiunque può accedere al tuo account
					utilizzando il tuo username
				</p>
			</div>
		</div>
	</div>
</template>
