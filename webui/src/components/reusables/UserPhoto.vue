<script setup>
import { ref, watchEffect } from "vue";
import fallbackUrl from "@/assets/images/default-avatar.png";

const props = defineProps({
	url: {
		type: String,
		default: null,
	},
	injectedFallbackUrl: {
		type: String,
		default: null,
	},
	size: {
		type: Number,
		default: 50, // default size in px
	},
});

const fb = props.injectedFallbackUrl || fallbackUrl;
const currentUrl = ref(fb);

watchEffect(() => {
	if (props.url) {
		// console.log(`${__API_URL__}/uploads${props.url}`);
		currentUrl.value = `${__API_URL__}/uploads${props.url}`;
	} else {
		currentUrl.value = fb;
	}
});

// handle image load error → fallback
const onError = () => {
	currentUrl.value = fb;
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
