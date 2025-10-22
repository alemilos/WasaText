<script setup>
import UserPhoto from "../../reusables/UserPhoto.vue";

const props = defineProps({
	users: {
		type: Object,
		default: [],
	},

	selectedUserId: {
		type: [String, Number, null],
		default: null,
	},
	onSelect: {
		type: Function,
		required: true,
	},
});
</script>

<template>
	<div class="usersSelectorList-container">
		<div
			v-for="user in users"
			:key="user.userId"
			class="userSelector-row-container"
			@click="props.onSelect(user.userId)"
			:style="{
				backgroundColor:
					user.userId === props.selectedUserId
						? 'var(--color-white-20)'
						: 'var(--color-white-5)',
			}"
		>
			<UserPhoto :url="user.photoPath" :size="36" />
			<span style="margin-left: 10px">{{ user.username }}</span>
		</div>
	</div>
</template>

<style scoped>
.usersSelector-container {
	display: flex;
	flex-direction: column;
	gap: 16px;
}

.usersSelectorList-container {
	display: flex;
	flex-direction: column;
	gap: 4px;
	height: 100%;
	overflow-y: scroll;
	max-height: 420px;
}
.userSelector-row-container {
	border-radius: 8px;
	display: flex;
	align-items: center;
	padding: 8px;
	cursor: pointer;
	background-color: var(--color-white-5);
}
.userSelector-row-container:hover {
	background-color: var(--color-white-20);
}
.userSelector-row-container span {
	color: var(--color-white);
}
</style>
