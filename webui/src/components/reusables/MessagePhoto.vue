<script setup>
import { ref, watchEffect } from "vue";
import fallbackUrl from "@/assets/images/default-message-image.png";

const props = defineProps({
	url: {
		type: String,
		default: null,
	},
	size: {
		type: Number,
		default: 50, // default size in px
	},
});

const currentUrl = ref(fallbackUrl);

watchEffect(() => {
	if (props.url) {
		console.log(`${__API_URL__}${props.url}`);
		currentUrl.value = `${__API_URL__}${props.url}`;
	} else {
		currentUrl.value = fallbackUrl;
	}
});

// handle image load error → fallback
const onError = () => {
	currentUrl.value = fallbackUrl;
};
</script>

<template>
	<div
		class="message-image-container"
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
.message-image-container {
	overflow: hidden;
	background-color: black;
}

.message-image-container img {
	width: 100%;
	height: 100%;
	object-fit: cover;
}
</style>
