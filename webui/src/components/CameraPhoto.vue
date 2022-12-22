<template>
	<div class="d-flex justify-content-center" v-if="streamingState === StreamingState.TURNED_OFF">
		<button @click="askCameraPermission" :disabled="loading" type="button" class="btn btn-primary btn-lg">
			Open camera
		</button>
	</div>

	<WarningMsg v-if="streamingState === StreamingState.ASKING_PERMISSION"
				msg="Please grant camera permission if requested!"/>

	<SuccessMsg v-if="streamingState === StreamingState.TURNING_ON" msg="Opening camera..."/>

	<div v-show="streamingState === StreamingState.TURNED_ON">
		<!-- Camera selector -->
		<div class="dropdown mb-2">
			<button class="btn btn-secondary dropdown-toggle" type="button" id="cameraSelectorButton"
					data-bs-toggle="dropdown" aria-expanded="false">
				{{ selectedCamera ? selectedCamera.label : 'No camera available' }}
			</button>
			<ul class="dropdown-menu" aria-labelledby="cameraSelectorButton">
				<li v-for="camera in availableCameras" :key="camera.deviceId" class="dropdown-item" @click="() => updateCamera(camera)">{{ camera.label }}</li>
			</ul>
		</div>

		<video ref="cameraStream" :style="{'aspect-ratio': aspectRatio}" class="video-stream" @canplay="startVideoPlaying">Video stream not available.</video>
		<div class="row-cols-3 mt-3">
			<button type="button" class="btn btn-primary m-2" @click="shootPhoto">Take photo</button>
			<button type="button" class="btn btn-outline-danger m-2" @click="goBack">Close camera</button>
		</div>
		<canvas ref="photoCanvas" style="display: none" width="500" :height="videoHeight"/>
	</div>

	<div v-if="streamingState === StreamingState.SHOT">
		<img :src="shotPhoto" class="video-stream" alt="Shot photo"/>
		<div class="row-cols-3 mt-3">
			<button :disabled="loading" type="button" class="btn btn-success m-2" @click="uploadPhoto">Upload photo
			</button>
			<button :disabled="loading" type="button" class="btn btn-outline-danger m-2" @click="askCameraPermission">
				Discard photo
			</button>
		</div>
	</div>

	<ErrorMsg :msg="errorMessage"/>
</template>

<script>
import ErrorMsg from "./ErrorMsg.vue";
import SuccessMsg from "./SuccessMsg.vue";
import WarningMsg from "./WarningMsg.vue";

const StreamingState = Object.freeze({
	TURNED_OFF: 0,
	ASKING_PERMISSION: 100,
	TURNING_ON: 200,
	TURNED_ON: 300,
	SHOT: 400,
});

export default {
	name: "CameraPhoto",
	emits: ['shooting', 'shot'],
	props: ['loading'],
	components: {WarningMsg, SuccessMsg, ErrorMsg},
	data() {
		return {
			errorMessage: null,
			StreamingState,
			videoHeight: 0,
			aspectRatio: 0,
			streamingState: StreamingState.TURNED_OFF,
			photoStream: null,
			shotPhoto: null,
			availableCameras: [],
			selectedCamera: null,
		};
	},
	methods: {
		async askCameraPermission() {
			this.streamingState = this.streamingState > StreamingState.ASKING_PERMISSION
				? StreamingState.TURNING_ON : StreamingState.ASKING_PERMISSION;
			this.$emit('shooting', true);

			try {
				if (!navigator.mediaDevices) {
					this.errorMessage = 'HTTPS (or localhost) is required to request camera permission';
				} else {
					await navigator.permissions.query({name: 'camera'});
					this.photoStream = await navigator.mediaDevices.getUserMedia({
						video: this.selectedCamera ? {deviceId: this.selectedCamera.deviceId} : true,
						audio: false,
					});
					this.streamingState = StreamingState.TURNING_ON;
					this.$refs.cameraStream.srcObject = this.photoStream;
					await this.$refs.cameraStream.play();

					// Now that we have camera permission, read other cameras available
					this.availableCameras = (await navigator.mediaDevices.enumerateDevices())
						.filter(camera => camera.kind === 'videoinput')
						.filter(camera => camera.label);
					if (!this.selectedCamera) this.selectedCamera = this.availableCameras[0];
				}
			} catch (err) {
				console.error(err);
				switch (err.name) {
					case 'NotAllowedError':
						this.errorMessage = 'Camera permission was not granted';
						break;
					case 'NotFoundError':
						this.errorMessage = 'No camera available';
						break;
					default:
						this.errorMessage = err.toString();
				}
			}
		},
		startVideoPlaying() {
			let videoHeight = this.$refs.cameraStream.videoHeight / (this.$refs.cameraStream.videoWidth / 500);

			// Firefox currently has a bug where the height can't be read from
			// the video, so we will make assumptions if this happens.

			if (isNaN(videoHeight)) {
				videoHeight = 500 / (4 / 3);
			}

			// Resize using aspect ratio
			this.aspectRatio = this.$refs.cameraStream.videoWidth / this.$refs.cameraStream.videoHeight;
			this.videoHeight = videoHeight;
			this.streamingState = StreamingState.TURNED_ON;
		},
		closeCamera() {
			if (this.photoStream) {
				this.photoStream.getTracks().forEach(track => track.stop());
			}
			this.photoStream = null;
			this.videoHeight = 0;
		},
		updateCamera(camera) {
			this.selectedCamera = camera;
			this.closeCamera();
			this.askCameraPermission();
		},
		shootPhoto() {
			const context = this.$refs.photoCanvas.getContext("2d");
			context.drawImage(this.$refs.cameraStream, 0, 0, 500, this.videoHeight);

			this.shotPhoto = this.$refs.photoCanvas.toDataURL('image/png');
			this.streamingState = StreamingState.SHOT;

			this.closeCamera();
		},
		uploadPhoto() {
			const file = dataURLtoBlob(this.shotPhoto);
			this.$emit('shot', [file]);
		},
		goBack() {
			this.closeCamera();
			this.$emit('shooting', false);
			this.streamingState = StreamingState.TURNED_OFF;
		},
	},
	unmounted() {
		this.closeCamera();
	},
};

function dataURLtoBlob(dataURL) {
	const byteString = atob(dataURL.split(',')[1]);
	const byteArray = new Uint8Array(byteString.length);
	for (let i = 0; i < byteString.length; i++) {
		byteArray[i] = byteString.charCodeAt(i);
	}
	return new Blob([byteArray], {type: 'image/png'});
}
</script>

<style scoped>
.video-stream {
	width: 500px;
	max-width: 90vw;
}
</style>
