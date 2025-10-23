<script setup>
import { ref, provide, onMounted, onUnmounted, markRaw } from "vue";
import ModalOverlay from "@/components/ui/ModalOverlay.vue";

const modal = ref(null);

function hidePageScroll() {
	document.body.classList.add("hide-scrollbar");
}

function showPageScroll() {
	document.body.classList.remove("hide-scrollbar");
}

function openModal(ModalComponent, props = {}) {
	hidePageScroll();
	modal.value = { ModalComponent: markRaw(ModalComponent), props };
}

function closeModal() {
	showPageScroll();
	modal.value = null;
}

// Optional: Close on Escape
function handleEscape(e) {
	if (e.code === "Escape") closeModal();
}

onMounted(() => window.addEventListener("keydown", handleEscape));
onUnmounted(() => window.removeEventListener("keydown", handleEscape));

// provide functions to all child components
provide("openModal", openModal);
provide("closeModal", closeModal);
provide("modal", modal);
</script>

<template>
	<slot />
	<ModalOverlay v-if="modal?.ModalComponent" @click="closeModal">
		<component :is="modal.ModalComponent" v-bind="modal.props" />
	</ModalOverlay>
</template>
