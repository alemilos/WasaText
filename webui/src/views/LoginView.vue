<script setup>
	import "./style/LoginView.css";
	import Logo from "../assets/images/logo.png";
	import { hackerText } from "../utils/text";
	import {useTemplateRef, onMounted} from 'vue'
	import {useRouter} from 'vue-router'
	import axios from '../services/axios'

	const title = useTemplateRef('title')
	const username = useTemplateRef('username')
	const router = useRouter()

	async function submit(e) {
			e.preventDefault();
			if (!username.value || username.value === "") {
				throw new Error("Invalid Username");
			}
			const res = await axios.post("/login", {username});
		
			if (res.status === 200){
				router.push("/")
			}
		}

	onMounted(() =>  {
		hackerText(title.value); // apply hacker animation
	})


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
