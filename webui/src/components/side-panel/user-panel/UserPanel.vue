<script setup>
import { computed } from "vue";
import CreateGroupIcon from "@/assets/icons/create-group.svg";
import CreateConvIcon from "@/assets/icons/create-conversation.svg";
import { usersStore } from "../../../stores/usersStore";
import { auth } from "../../../stores/authStore";
import UserPhoto from "../../reusables/UserPhoto.vue";

const userId = computed(() => auth.userId);
const profilePhotoPath = computed(() => usersStore.getPhotoUrl(userId.value));
</script>

<template>
	<div class="userPanel-container">
		<button>
			<img :src="CreateConvIcon" />
		</button>
		<button>
			<img :src="CreateGroupIcon" />
		</button>

		<button class="userPanel-profile-btn">
			<UserPhoto :url="profilePhotoPath" :size="42" />
		</button>
	</div>
</template>

<style scoped>
.userPanel-container {
	display: flex;
	align-items: center;
	padding: 16px 20px;
	gap: 12px;
	border-top: 1px solid var(--color-white-20);
}

.userPanel-container button {
	width: 42px;
	height: 42px;
	min-width: 42px;
	min-height: 42px;
	border-radius: 50%;
	background-color: var(--color-green-primary);
	border: none;
	display: flex;
	align-items: center;
	justify-content: center;
}

.userPanel-container button:hover {
	opacity: 0.7;
}

.userPanel-profile-btn {
	margin-left: auto;
	background-color: black !important; /* Override green bg */
}

.userPanel-profile-btn:hover {
	opacity: 0.7;
}
</style>
