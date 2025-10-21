<script setup></script>

<template>
	<div class="unselectedConversation-container">
		<div class="unselectedConversation-message">
			<p>
				Seleziona una conversazione o inizia a chattare per dare vita
				alla tua conversazione!
			</p>
		</div>
	</div>
</template>

<style scoped>
.unselectedConversation-container {
	width: 100%;
	height: 100%;
	display: flex;
	justify-content: center;
	align-items: end;
}

.unselectedConversation-message {
	text-align: center;
	background-color: rgba(0, 0, 0, 0.9);
	padding: 30px 30px 14px 30px;
	border-radius: 12px;
	box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
	font-size: 16px;
	color: var(--color-green-primary);
	border: 1px solid var(--color-white-20);
}
</style>
