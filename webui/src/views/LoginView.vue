<script setup>
import "./style/LoginView.css";
import Logo from "../assets/images/logo.png";
import { hackerText } from "../utils/text";
import axios from "../services/axios";
</script>

<script>
export default {
	data() {
		return {
			logo: Logo,
		};
	},

	methods: {
		submit(e) {
			e.preventDefault();
			const username = this.$refs.username.value;
			if (!username || username === "") {
				throw new Error("Invalid Username");
			}

			axios.post("/login", {});
		},
	},

	mounted() {
		hackerText(this.$refs.title); // apply hacker animation
	},
};
</script>

<template>
	<div class="login-container">
		<div class="form-container">
			<div class="form-heading">
				<img :src="logo" alt="logo" />
				<h1 ref="title">Accedi a WasaText</h1>
				<p>
					Chatta in modo <span>sicuro</span> e <span>veloce</span>,
					con amici e familiari.
				</p>
			</div>
			<form>
				<input type="text" placeholder="username" ref="username" />
				<button @click="submit">Accedi</button>
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
