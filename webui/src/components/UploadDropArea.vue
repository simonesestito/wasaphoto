<script>
import UploadDropIcon from "./UploadDropIcon.vue";

export default {
	name: "UploadDropArea",
	components: {UploadDropIcon},
	props: ['disabled'],
	emits: ['drop'],
	data() {
		return {
			droppingStatus: 0,
		};
	},
	methods: {
		setActive() {
			if (!this.disabled) this.droppingStatus++;
		},
		setInactive() {
			if (!this.disabled) {
				setTimeout(() => this.droppingStatus--, 200);
			}
		},
		onDrop(event) {
			if (this.disabled) return;

			event.preventDefault();
			this.droppingStatus = false;

			// Compatibility with both DataTransfer and DataTransferItemList
			let files;
			if (event.dataTransfer.items) {
				// Use DataTransferItemList interface to access the file(s)
				files = [...event.dataTransfer.items]
					.filter(item => item.kind === 'file')
					.map(item => item.getAsFile());
			} else {
				// Use DataTransfer interface to access the file(s)
				files = [...event.dataTransfer.files];
			}

			this.$emit('drop', files);
		},
	},
}
</script>

<template>
	<div class="drag-area"
		 @dragenter.prevent="setActive"
		 @dragleave.prevent="setInactive"
		 @dragover.prevent=""
		 @drop.prevent="onDrop">

		<slot @dragenter.prevent="setActive"
			  @dragleave.prevent="setInactive"
			  @dragover.prevent="setActive"
			  @drop.prevent="onDrop"/>

		<div :class="{active: droppingStatus}" class="drag-area-indicator">
			<UploadDropIcon/>
		</div>
	</div>
</template>

<style scoped>
.drag-area {
	width: 100%;
	height: 90vh;
}

.drag-area-indicator {
	position: absolute;
	top: 0;
	left: 0;
	width: 100vw;
	height: 100vh;

	display: none;
	opacity: 0;
	transition: opacity 0.15s ease-in;

	background-color: rgba(48, 202, 214, 0.4);
}

.drag-area-indicator.active {
	display: block;
	opacity: 1;
}
</style>
