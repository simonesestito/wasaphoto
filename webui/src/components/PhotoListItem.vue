<script>
import {LikesService} from "../services";
import router from "../router";
import {toRefs} from "vue";
import ErrorMsg from "./ErrorMsg.vue";
import LoadingSpinner from "./LoadingSpinner.vue";

export default {
	name: "PhotoListItem",
	components: {LoadingSpinner, ErrorMsg},
	props: ['photo', 'showAuthor'],
	emits: ['error'],
	data() {
		return {
			loading: false,
		};
	},
	setup(props) {
		return toRefs(props);
	},
	methods: {
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
		async goToAuthor() {
			await router.push(`/users/${this.photo.author.id}`);
		}
	}
}
</script>

<template>
	<div class="p-4 mt-3">
		<div v-if="photo" class="col">
			<div v-if="showAuthor" class="row user-header" data-bs-dismiss="modal" @click="goToAuthor">
				<p>{{photo.author.name}} {{photo.author.surname}} (@{{photo.author.username}})</p>
			</div>
			<div class="photo-content">
				<img :src="photo.imageUrl" alt="User photo">
			</div>
			<div class="row actions-row">
				<p class="likes" @click="doLike(!photo.liked)">
					<svg class="feather" :class="{ active: photo.liked, disabled: loading }"><use href="/feather-sprite-v4.29.0.svg#thumbs-up"/></svg>
					<span>{{ photo.likesCount }}</span>
				</p>

				<p class="comments" @click="goToComments" data-bs-dismiss="modal">
					<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#message-square"/></svg>
					<span>{{ photo.commentsCount }}</span>
				</p>
			</div>
		</div>
	</div>
</template>

<style scoped>
.user-header {
	cursor: pointer;
	font-size: 1.2rem;
}

.photo-content {
	max-width: 300px;
	height: 300px;
	margin: 0 auto;
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

.user-header {
	height: 40px;
}

.actions-row {
	display: flex;
	justify-content: center;
}
</style>

