<script>
import UploadDropIcon from "./UploadDropIcon.vue";
export default {
	name: "UploadDropArea",
	components: {UploadDropIcon},
	emits: ['drop'],
	data() {
		return {
			droppingStatus: false,
		};
	},
	methods: {
		setActive() {
			this.droppingStatus = true;
		},
		setInactive() {
			this.droppingStatus = false;
		},
		onDrop(event) {
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
		 :class="{active: droppingStatus}"
		 @dragenter.prevent="setActive"
		 @dragleave.prevent="setInactive"
		 @dragover.prevent="setActive"
		 @drop.prevent="onDrop">
		<UploadDropIcon />
	</div>
</template>

<style scoped>
.drag-area {
	position: absolute;
	top: 0;
	left: 0;
	width: 100vw;
	height: 100vh;
	background-color: rgba(48,202,214,0.4);
	opacity: 0;
	transition: opacity 0.15s ease-in;
	z-index: -1;
}

.drag-area.active {
	opacity: 1;
	z-index: 10;
}
</style>
