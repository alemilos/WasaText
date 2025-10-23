<script setup>
const props = defineProps({
	label: String,
	type: {
		type: String,
		default: "text",
	},
	isFloat: Boolean,
	value: [String, Number],
	placeholder: String,
	readOnly: {
		type: Boolean,
		default: false,
	},
	inputClassName: String,
	labelClassName: String,
});

const emit = defineEmits(["update:value"]);

function handleInput(event) {
	emit("update:value", event.target.value);
}
</script>

<template>
	<div class="labeledInput">
		<span class="labeledInput-label" :class="labelClassName">{{
			label
		}}</span>

		<input
			:type="type"
			:step="isFloat ? 'any' : null"
			:value="value"
			:readonly="readOnly"
			:placeholder="placeholder"
			class="labeledInput-input"
			:class="inputClassName"
			@input="handleInput"
		/>
	</div>
</template>

<style scoped>
.labeledInput {
	display: flex;
	align-items: center;
	gap: 8px;
}

.labeledInput-label {
	color: white;
}

.labeledInput-input {
	width: 100%;
	height: 100%;
	padding: 8px 12px;
	border: none;
	border-radius: 16px;
	background-color: #1a1a1a; /* bgPrimary equivalent */
	color: white; /* textWhite equivalent */
	outline: none;
	transition: background-color 0.2s ease;
}

.labeledInput-input:hover {
	background-color: rgba(255, 255, 255, 0.1);
}

.labeledInput-input:focus {
	background-color: rgba(255, 255, 255, 0.15);
}

.labeledInput-input[readonly] {
	cursor: default;
	background-color: #1a1a1a;
}

.labeledInput-input::placeholder {
	color: rgba(255, 255, 255, 0.4); /* textWhiteDisabled equivalent */
}
</style>
