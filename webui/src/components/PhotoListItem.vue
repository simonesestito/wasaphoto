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
					await LikesService.likePhoto(this.photo.id);
					this.photo.likesCount++;
				} else {
					await LikesService.unlikePhoto(this.photo.id);
					this.photo.likesCount--;
				}

				this.photo.liked = like;
			} catch (err) {
				this.$emit('error', err);
			} finally {
				this.loading = false;
			}
		},
		async goToComments() {
			await router.push(`/photos/${this.photo.id}/comments`);
		},
		async deletePhoto() {
			this.loading = true;
			this.$emit('error', '');
			try {
				await PhotosService.deletePhoto(this.photo.id);
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
			window.open(this.photo.imageUrl, '_blank').focus();
		},
	},
	computed: {
		isMine() {
			if (!this.photo) return false;
			return getCurrentUID() === this.photo.author.id;
		}
	}
}
</script>

<template>
	<div class="p-4 mt-3" :class="{card: showAuthor}">
		<div style="display: none" data-bs-dismiss="modal" ref="close"/> <!-- Close hidden HTML element -->
		<div v-if="photo" class="col">
			<UserNameHeader v-if="showAuthor" :user="photo.author"/>
			<div class="photo-content" @click="openImageNewTab">
				<img :src="photo.imageUrl" alt="User photo" loading="lazy">
			</div>
			<p class="post-date">{{ formatDate(photo.publishDate) }}</p>
			<div class="row actions-row">
				<p class="likes" @click="doLike(!photo.liked)">
					<svg class="feather" :class="{ active: photo.liked, disabled: loading }">
						<use href="/feather-sprite-v4.29.0.svg#thumbs-up"/>
					</svg>
					<span>{{ photo.likesCount }}</span>
				</p>

				<p class="comments" @click="goToComments" data-bs-dismiss="modal">
					<svg class="feather">
						<use href="/feather-sprite-v4.29.0.svg#message-square"/>
					</svg>
					<span>{{ photo.commentsCount }}</span>
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

