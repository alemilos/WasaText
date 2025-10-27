<script setup>
import { computed } from "vue";
import Modal from "../ui/Modal.vue";
import ModalTitle from "../ui/ModalTitle.vue";
import UserPhoto from "../reusables/UserPhoto.vue";
import { usersStore } from "../../stores/usersStore";
import { formatDate } from "../../utils/time";

const props = defineProps({
	reactions: {
		type: Array,
		required: true,
	},
});

const detailedReactions = computed(() =>
	props.reactions.map((r) => ({
		...r,
		username: usersStore.getUsername(r.authorId),
		photo: usersStore.getPhotoUrl(r.authorId),
		formattedTime: formatDate(r.createdAt),
	}))
);
</script>

<template>
	<Modal>
		<div class="reactionsModal-container">
			<ModalTitle text="Reazioni al messaggio" />

			<div v-if="!reactions.length" class="no-reactions">Nessuna reazione al momento.</div>

			<div v-else class="reactions-list">
				<div
					v-for="reaction in detailedReactions"
					:key="reaction.authorId + reaction.createdAt"
					class="reaction-row"
				>
					<UserPhoto :url="reaction.photo" :size="32" />
					<div class="reaction-info">
						<span class="username">{{ reaction.username }}</span>
						<span class="time">{{ reaction.formattedTime }}</span>
					</div>
					<span class="emoji">{{ reaction.emoji }}</span>
				</div>
			</div>
		</div>
	</Modal>
</template>

<style scoped>
.reactionsModal-container {
	padding: 16px;
	color: var(--color-white);
}

.reactions-list {
	display: flex;
	flex-direction: column;
	gap: 10px;
	margin-top: 12px;
	max-height: 500px;
	overflow-y: scroll;
}

.reaction-row {
	display: flex;
	align-items: center;
	justify-content: space-between;
	background: rgba(255, 255, 255, 0.05);
	padding: 8px 12px;
	border-radius: 10px;
	transition: background 0.2s;
}

.reaction-row:hover {
	background: rgba(255, 255, 255, 0.1);
}

.reaction-info {
	display: flex;
	flex-direction: column;
	flex-grow: 1;
	margin-left: 8px;
}

.username {
	font-weight: 600;
	font-size: 14px;
}

.time {
	font-size: 12px;
	opacity: 0.7;
}

.emoji {
	font-size: 20px;
	margin-left: 12px;
}

.no-reactions {
	text-align: center;
	padding: 20px;
	font-size: 14px;
	color: var(--color-red);
}
</style>
