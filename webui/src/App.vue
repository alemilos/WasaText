<script setup>
import { RouterView } from "vue-router";
import { watch, onUnmounted } from "vue";
import { auth } from "./stores/authStore";
import { usersStore } from "./stores/usersStore";
import ModalProvider from "./providers/ModalProvider.vue";
import { conversationsStore } from "./stores/conversationsStore";

let intervalId;

watch(
	() => auth.userId,
	async (newUserId) => {
		if (newUserId) {
			await Promise.all([
				usersStore.loadUsers(),
				conversationsStore.loadConversations(),
			]);

			// clear prev intervals
			if (intervalId) clearInterval(intervalId);

			intervalId = setInterval(async () => {
				await Promise.all([
					usersStore.loadUsers(),
					conversationsStore.loadConversations(),
				]);
			}, 3000);
		} else {
			// on logout clear state
			if (intervalId) {
				clearInterval(intervalId);
				intervalId = null;
			}
		}
	},

	{ immediate: true }
);

// Clean up on component unmount
onUnmounted(() => {
	if (intervalId) clearInterval(intervalId);
});
</script>

<script>
export default {};
</script>

<template>
	<ModalProvider>
		<RouterView />
	</ModalProvider>
</template>
