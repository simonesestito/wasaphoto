<!-- Heavily edited, from https://codepen.io/megaroeny/pen/GoyYQG -->

<template>
	<div class="upload-btn-container">

		<svg class="upload-box" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 17"
			 xml:space="preserve"
			 style="clip-rule:evenodd;fill-rule:evenodd;stroke-linejoin:round;stroke-miterlimit:1.41">
			<path id="upload-box"
				  d="M26 3.82l3.16 1.18 -13.16 4.91 -13.16-4.91 3.16-1.18 0-3.07 -6 2.25 0 8 16 6 16-6 0-8 -6-2.25 0 3.07Z"
				  fill="#ffffff"/>
		</svg>

		<svg class="upload-arrow" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"
			 xml:space="preserve"
			 style="clip-rule:evenodd;fill-rule:evenodd;stroke-linejoin:round;stroke-miterlimit:1.41">
			<path id="upload-arrow" d="M6 16l4 0 0-8 6 0 -8-8 -8 8 6 0 0 8Z" fill="#ffffff"/>
		</svg>

	</div>
</template>

<script>
export default {
	name: "UploadDropIcon"
}
</script>

<style scoped>

@keyframes uploading-arrow {
	0% {
		padding-bottom: 40px;
		opacity: 0;
	}
	20% {
		opacity: 1;
	}
	50% {
		opacity: 1;
	}
	75% {
		opacity: 0;
	}
	100% {
		padding-bottom: 20vh;
	}
}
.upload-btn-container {
	overflow: hidden;
	display: block;
	width: 100px;
	height: 100px;
}
.upload-btn-container svg {
	position: absolute;
	top: 0;
	bottom: 0;
	right: 0;
	left: 0;
	margin: auto;
	width: 100px;
}
.upload-btn-container svg.upload-box {
	transform: scale(150%);
}
.upload-btn-container svg.upload-arrow {
	opacity: 0;
	-webkit-animation: uploading-arrow 1.5s ease-out infinite;
	animation: uploading-arrow 1.5s ease-out infinite;
}
</style>
