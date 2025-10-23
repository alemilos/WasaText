<script setup>
import ChatPanel from "../components/chat-panel/ChatPanel.vue";
import SidePanel from "../components/side-panel/SidePanel.vue";
import Logo from "../assets/images/logo.png";
import "./style/HomeView.css";
import { createUsers } from "../services/helpers";
import { useToast } from "vue-toast-notification";

const $toast = useToast();

async function handleCreateUsers() {
	await createUsers();
	$toast.success("10 utenti creati, alcuni potrebbero già esistere");
}
</script>

<template>
	<div class="home-container">
		<div class="home-heading">
			<img class="img-logo" :src="Logo" />
			<div class="helpers-container">
				<button class="helper-button" @click="handleCreateUsers">
					Crea utenti velocemente
				</button>
			</div>
		</div>
		<div class="main-container">
			<SidePanel />
			<ChatPanel />
		</div>
	</div>
</template>
<style scoped>
.home-heading {
	display: flex;
	align-items: center;
}
.helpers-container {
	margin-left: auto;
	display: flex;
	gap: 8px;
	align-items: center;
	color: var(--color-white);
}
.helpers-container span {
	margin-top: 12px;
	font-size: 18px;
	margin-right: 20px;
}

.helpers-container button {
	border: none;
	border-radius: 8px;
	padding: 6px 16px;
	margin-top: 16px;
	margin-left: auto;
}

.helpers-container button:hover {
	opacity: 0.7;
}
.helper-button {
	background-color: var(--color-green-primary);
}
.helper-button:hover {
	opacity: 0.7;
}
.helper-button:disabled {
	opacity: 0.4;
	cursor: not-allowed;
}
</style>
