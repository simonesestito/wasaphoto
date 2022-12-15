<script>
import PageSkeleton from "../components/PageSkeleton.vue";
import {PhotosService, UsersService} from "../services";
import LoadingSpinner from "../components/LoadingSpinner.vue";
import SuccessMsg from "../components/SuccessMsg.vue";
import ErrorMsg from "../components/ErrorMsg.vue";
import UploadDropArea from "../components/UploadDropArea.vue";
import router from "../router";
import {getCurrentUID} from "../services/auth-store";
import CameraPhoto from "../components/CameraPhoto.vue";

export default {
	name: "PhotoUploadView",
	components: {CameraPhoto, UploadDropArea, LoadingSpinner, PageSkeleton, SuccessMsg, ErrorMsg},
	data() {
		return {
			loading: false,
			success: false,
			errorMessage: null,
			isShooting: false,
			uploadProgress: 0,
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
				} else if (fileList[0].size > 20 * 1024 * 1024) {
					this.errorMessage = 'File is too large (max allowed is 20MB)';
				} else {
					await PhotosService.uploadPhoto(fileList[0], progress => this.uploadProgress = progress);
					this.success = true;
					if (this.$refs.camera) this.$refs.camera.goBack();
				}
			} catch (err) {
				this.errorMessage = err.toString();
			} finally {
				this.loading = false;
				if (this.$refs.input) this.$refs.input.value = '';
			}
		},
		async goToMyProfile() {
			const myProfile = await UsersService.getUserProfile(getCurrentUID());
			await router.push(`/users/${myProfile.username}`);
		}
	},
}
</script>

<template>
	<UploadDropArea @drop="uploadPhoto" :disabled="isShooting">
		<PageSkeleton title="Upload New Photo" :actions="[{text:'My profile', onClick: goToMyProfile}]">
			<ErrorMsg :msg="errorMessage"/>
			<SuccessMsg v-if="success" msg="Upload succeeded"/>

			<div v-if="loading">
				<LoadingSpinner />
				<p class="progress-status">{{ uploadProgress < 100 ? uploadProgress.toFixed(0) + '%' : 'Processing...' }}</p>
			</div>

			<div v-if="!isShooting">
				<input type="file" accept="image/*" ref="input" class="photo-upload-input" @change="onFilePicked">

				<div class="photo-upload-box" :class="{disabled: loading}"
					 @click="requestPhotoFile">
					Click to pick an image or drop it here
				</div>

				<p class="or-option-text"><i>Or...</i></p>
			</div>

			<CameraPhoto ref="camera" :loading="loading" @shooting="e => isShooting = e" @shot="uploadPhoto"/>
		</PageSkeleton>
	</UploadDropArea>
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

.or-option-text {
	text-align: center;
}

.progress-status {
	text-align: center;
	font-size: 1.2rem;
}
</style>
