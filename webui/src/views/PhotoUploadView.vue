<script>
import PageSkeleton from "../components/PageSkeleton.vue";
import {PhotosService} from "../services";
import LoadingSpinner from "../components/LoadingSpinner.vue";
import SuccessMsg from "../components/SuccessMsg.vue";
import ErrorMsg from "../components/ErrorMsg.vue";
import UploadDropArea from "../components/UploadDropArea.vue";

export default {
	name: "PhotoUploadView",
	components: {UploadDropArea, LoadingSpinner, PageSkeleton, SuccessMsg, ErrorMsg},
	data() {
		return {
			loading: false,
			success: false,
			errorMessage: null
		};
	},
	methods: {
		requestPhotoFile() {
			if (!this.loading)
				this.$refs.input.click();
		},
		onFilePicked() {
			this.uploadPhoto(this.$refs.input.files);
		},
		async uploadPhoto(fileList) {
			if (this.loading) return;

			this.loading = true;
			this.errorMessage = null;
			this.success = false;

			try {
				if (!fileList || fileList.length === 0) {
					this.errorMessage = 'No photo found to upload';
				} else if (fileList.length > 1) {
					this.errorMessage = 'More than one photo selected';
					this.$refs.input.value = '';
				} else {
					await PhotosService.uploadPhoto(fileList[0]);
					this.success = true;
				}
			} catch (err) {
				this.errorMessage = err.toString();
			} finally {
				this.loading = false;
				this.$refs.input.value = '';
			}
		},
	},
}
</script>

<template>
	<UploadDropArea @drop="uploadPhoto" />

	<PageSkeleton title="Upload New Photo">
		<ErrorMsg :msg="errorMessage"/>
		<SuccessMsg v-if="success" msg="Upload succeeded"/>
		<LoadingSpinner v-if="loading"/>

		<input type="file" accept="image/*" ref="input" class="photo-upload-input" @change="onFilePicked">

		<div class="photo-upload-box" :class="{disabled: loading}"
			 @click="requestPhotoFile">
			Click to pick an image or drop it here
		</div>
	</PageSkeleton>
</template>

<style scoped>
.photo-upload-input {
	display: none;
}

.photo-upload-box {
	width: 80%;
	height: 200px;
	line-height: 200px;
	margin: 8px auto;
	background-color: #eeeeee;
	border: 1px solid #7a7a7a;
	border-radius: 8px;
	text-align: center;
	cursor: pointer;
}

.photo-upload-box.disabled {
	cursor: no-drop;
}
</style>
