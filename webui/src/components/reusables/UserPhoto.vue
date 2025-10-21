<script setup>
import { computed, ref, watchEffect } from "vue";
import userImage from "@/assets/images/default-avatar.png";
import groupImage from "@/assets/images/default-group.png";

const props = defineProps({
	url: {
		type: String,
		default: null,
	},
	type: {
		type: String,
		default: "private",
	},
	size: {
		type: Number,
		default: 50, // default size in px
	},
});

const fallbackImage = computed(() =>
	props.type === "group" ? groupImage : userImage
);
const currentUrl = ref(fallbackImage.value);

watchEffect(() => {
	if (props.url) {
		// console.log(`${__API_URL__}/uploads${props.url}`);
		currentUrl.value = `${__API_URL__}/uploads${props.url}`;
	} else {
		currentUrl.value = fallbackImage.value;
	}
});

// handle image load error → fallback
const onError = () => {
	currentUrl.value = fallbackImage.value;
};
</script>

<template>
	<div
		class="user-avatar"
		:style="{
			width: size + 'px',
			height: size + 'px',
			minWidth: size + 'px',
			minHeight: size + 'px',
		}"
	>
		<img :src="currentUrl" @error="onError" />
	</div>
</template>

<style scoped>
.user-avatar {
	border-radius: 50%;
	overflow: hidden;
	background-color: black;
}

.user-avatar img {
	width: 100%;
	height: 100%;
	object-fit: cover;
}
</style>
