<script>
import {LikesService, PhotosService} from "../services";
import router from "../router";
import {toRefs} from "vue";
import {getCurrentUID} from "../services/auth-store";
import {formatDate} from "../services/format-date";
import UserNameHeader from "./UserNameHeader.vue";

export default {
	name: "PhotoListItem",
	components: {UserNameHeader},
	props: ['photo', 'showAuthor'],
	emits: ['error', 'refresh'],
	data() {
		return {
			photoData: this.photo, // Use another internal variable for changes
			loading: false,
		};
	},
	setup(props) {
		return toRefs(props);
	},
	methods: {
		formatDate(date) {
			return formatDate(date);
		},
		async doLike(like) {
			if (this.loading) return;

			this.loading = true;
			this.$emit('error', '');
			try {
				if (like) {
					await LikesService.likePhoto(this.photoData.id);
					this.photoData.likesCount++;
				} else {
					await LikesService.unlikePhoto(this.photoData.id);
					this.photoData.likesCount--;
				}

				this.photoData.liked = like;
			} catch (err) {
				this.$emit('error', err);
			} finally {
				this.loading = false;
			}
		},
		async goToComments() {
			await router.push(`/photos/${this.photoData.id}/comments`);
		},
		async deletePhoto() {
			this.loading = true;
			this.$emit('error', '');
			try {
				await PhotosService.deletePhoto(this.photoData.id);
				this.$emit('refresh');
			} catch (err) {
				this.$emit('error', err.toString());
			} finally {
				this.close();
				this.loading = false;
			}
		},
		close() {
			// Check if we're inside a modal
			if (this.$el.closest('.modal') !== null) {
				// Dismiss dialog
				this.$refs.close.click();
			} else {
				// Pop router stack
				router.back();
			}
		},
		openImageNewTab() {
			window.open(this.photoData.imageUrl, '_blank').focus();
		},
	},
	computed: {
		isMine() {
			if (!this.photoData) return false;
			return getCurrentUID() === this.photoData.author.id;
		}
	}
}
</script>

<template>
	<div class="p-4 mt-3" :class="{card: showAuthor}">
		<div style="display: none" data-bs-dismiss="modal" ref="close"/> <!-- Close hidden HTML element -->
		<div v-if="photoData" class="col">
			<UserNameHeader v-if="showAuthor" :user="photoData.author"/>
			<div class="photo-content" @click="openImageNewTab">
				<img :src="photoData.imageUrl" alt="User photo" loading="lazy">
			</div>
			<p class="post-date">{{ formatDate(photoData.publishDate) }}</p>
			<div class="row actions-row">
				<p class="likes" @click="doLike(!photoData.liked)">
					<svg class="feather" :class="{ active: photoData.liked, disabled: loading }">
						<use href="/feather-sprite-v4.29.0.svg#thumbs-up"/>
					</svg>
					<span>{{ photoData.likesCount }}</span>
				</p>

				<p class="comments" @click="goToComments" data-bs-dismiss="modal">
					<svg class="feather">
						<use href="/feather-sprite-v4.29.0.svg#message-square"/>
					</svg>
					<span>{{ photoData.commentsCount }}</span>
				</p>
			</div>
			<div class="row d-flex justify-content-end">
				<button v-if="isMine" @click="deletePhoto" type="button" :disabled="loading" class="btn btn-outline-danger">Delete</button>
			</div>
		</div>
	</div>
</template>

<style scoped>
.card {
	max-width: 350px;
	margin: 0 auto;
	border: 1px gray solid;
}

.photo-content {
	max-width: 300px;
	height: 300px;
	margin: 0 auto;
	position: relative;
}

.post-date {
	text-align: center;
}

.photo-content img {
	max-width: 300px;
	max-height: 300px;
	position: absolute;
	margin: auto;
	top: 0;
	left: 0;
	right: 0;
	bottom: 0;
	cursor: zoom-in;
}

.likes, .comments {
	align-items: center;
	display: inline-flex;
	flex-wrap: wrap;
	width: auto;
	margin-right: 16px;
	margin-top: 8px;
}

.likes .feather, .comments .feather {
	display: inline-block;
	cursor: pointer;
	width: 1.5rem;
	height: 1.5rem;
	margin-right: 10px;
}

.feather.active {
	fill: #1976d2;
}

.feather.disabled {
	fill: lightgray;
	opacity: 0.8;
	cursor: not-allowed !important;
}

.comments {
	cursor: pointer;
}

.actions-row {
	display: flex;
	justify-content: center;
}
</style>

