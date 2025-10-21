<script setup>
import { RouterView } from "vue-router";
import { watch } from "vue";
import { auth } from "./stores/authStore";
import { usersStore } from "./stores/usersStore";

watch(
	() => auth.userId,
	async (newUserId) => {
		if (newUserId) {
			await usersStore.loadUsers();
		}
	},
	{ immediate: true }
);
</script>

<script>
export default {};
</script>

<template>
	<RouterView />
</template>

<!-- await usersStore.loadUsers(); -->
