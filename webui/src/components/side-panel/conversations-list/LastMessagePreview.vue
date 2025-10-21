<script setup>
import "./style/Conversation.css";
import imgIcon from "@/assets/images/image-icon.png";
import { formatDate } from "@/utils/time";
import { computed } from "vue";

const props = defineProps({
	conversationName: { type: String, default: "Sconosciuto" }, // group name | username
	type: { type: String, default: "text" }, // text | image
	content: { type: String, default: "" },
	createdAt: { type: String, default: null },
	createdBy: { type: Number, default: null },
});

const formattedDate = computed(() => {
	return props.createdAt ? formatDate(props.createdAt) : "";
});
</script>

<template>
	<div class="lastMessagePreview-container">
		<span class="lastMessagePreview-name">{{ conversationName }}</span>

		<!-- message preview -->
		<div class="lastMessagePreview-content">
			<template v-if="type === 'text'">
				<p
					v-if="content && content.trim() !== ''"
					class="lastMessagePreview-text-content"
				>
					{{ content }}
				</p>
				<p v-else class="lastMessagePreview-text-empty">
					La chat è vuota
				</p>
			</template>
			<template v-else-if="type === 'image'">
				<div class="lastMessagePreview-image-container">
					<img
						:src="imgIcon"
						alt=""
						class="lastMessagePreview-image-icon"
					/>
					<p>Foto</p>
				</div>
			</template>
		</div>

		<!-- date -->
		<span class="lastMessagePreview-date">{{ formattedDate }}</span>
	</div>
</template>
