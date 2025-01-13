<script setup>
	import ChatPanel from '../components/home/ChatPanel.vue';
	import SidePanel from '../components/home/SidePanel.vue';
	import Logo from '../assets/images/logo.png'
	import "./style/HomeView.css"
</script>

<template>
	<div class="home-container">
		<img class="img-logo" :src="Logo"  />
		<div class="main-container">
			<SidePanel/>
			<ChatPanel/>
		</div>
	</div>
</template>

