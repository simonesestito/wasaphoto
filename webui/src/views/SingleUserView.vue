<script>
import {FollowService, PhotosService, UsersService} from "../services";
import LoadingSpinner from "../components/LoadingSpinner.vue";
import ErrorMsg from "../components/ErrorMsg.vue";
import PageSkeleton from "../components/PageSkeleton.vue";
import {BanService} from "../services/ban";
import ShowMore from "../components/ShowMore.vue";
import PhotoListItem from "../components/PhotoListItem.vue";

export default {
	name: 'SingleUserView',
	components: {PhotoListItem, ShowMore, PageSkeleton, ErrorMsg, LoadingSpinner},
	props: ['username'],
	data() {
		return {
			loading: false,
			errorMessage: null,
			user: null,
			photos: [],
			photosCursor: null,
		};
	},
	methods: {
		async refresh(username) {
			// Set default value
			if (!username)
				username = this.$route.params.username;

			// Handle double loading or null (/undefined) parameters.
			if (!username /* with default value */ || this.loading) {
				return;
			}

			this.loading = true;
			this.errorMessage = null;
			this.user = null;
			this.photos = [];

			try {
				this.user = await UsersService.getByUsername(username);
				if (this.user == null) {
					this.errorMessage = 'User not found';
				} else {
					this.photosCursor = null;
					await this.loadMore();
				}
			} catch (err) {
				this.errorMessage = err.toString();
			} finally {
				this.loading = false;
			}
		},
		async setFollow(follow) {
			if (this.loading) return;
			this.loading = true;
			this.errorMessage = null;
			try {
				if (follow) {
					await FollowService.followUser(this.user.id);
				} else {
					await FollowService.unfollowUser(this.user.id);
				}

				await this.refresh(this.user.username);
			} catch (err) {
				this.errorMessage = err.toString();
			} finally {
				this.loading = false;
			}
		},
		async setBan(ban) {
			if (this.loading) return;
			this.loading = true;
			this.errorMessage = null;
			try {
				if (ban) {
					await BanService.banUser(this.user.id);
				} else {
					await BanService.unbanUser(this.user.id);
				}

				this.user.banned = ban;
			} catch (err) {
				this.errorMessage = err.toString();
			} finally {
				this.loading = false;
			}
		},
		async loadMore() {
			this.loading = true;
			this.errorMessage = null;

			try {
				const photoResponse = await PhotosService.listUserPhotos(this.user.id, this.photosCursor);
				this.photosCursor = photoResponse.nextPageCursor;
				this.photos.push(...photoResponse.pageData);
			} catch (err) {
				this.errorMessage = err.toString();
			} finally {
				this.loading = false;
			}
		},
	},
	computed: {
		displayUserName() {
			if (this.user) {
				return `${this.user.name} ${this.user.surname} (@${this.user.username})`
			} else {
				return '';
			}
		},
		followButton() {
			if (!this.user) {
				return null;
			} else if (this.user.following) {
				return {text: 'Unfollow', onClick: () => this.setFollow(false)};
			} else {
				return {text: 'Follow', onClick: () => this.setFollow(true)};
			}
		},
		banButton() {
			if (!this.user) {
				return null;
			} else if (this.user.banned) {
				return {text: 'Unban', onClick: async () => this.setBan(false)};
			} else {
				return {text: 'Ban', onClick: async () => this.setBan(true)};
			}
		},
	},
	mounted() {
		this.refresh();
	},
	created() {
		this.$watch(
			() => this.$route.params,
			(params, _) => {
				// react to route changes...
				this.refresh(params.username);
			}
		)
	},
}
</script>

<template>
	<PageSkeleton :title="displayUserName" :main-action="followButton" :actions="banButton ? [banButton] : []">
		<LoadingSpinner v-if="loading"/>
		<ErrorMsg :msg="errorMessage"/>

		<!-- User profile -->
		<div v-if="user">
			<div class="row-cols-md-3 d-flex justify-content-evenly">
				<span>
					<RouterLink :to="`/users/${this.user.username}/followers`"><b>Followers: </b> {{
							this.user.followersCount
						}}
					</RouterLink>
				</span>
				<span>
					<RouterLink
						:to="`/users/${this.user.username}/followings`"><b>Followings: </b> {{
							this.user.followingsCount
						}}
					</RouterLink>
				</span>
				<span><b>Posts count: </b> {{ this.user.postsCount }}</span>
			</div>

			<!-- Photos list -->
			<div>
				<div class="posts-grid mt-3">
					<div v-for="photo in photos"
						 :key="photo.id"
						 class="posts-grid-item"
						 data-bs-toggle="modal"
						 :data-bs-target="`#photo-modal-${photo.id}`"
						 :style="{backgroundImage: `url(${photo.imageUrl})`}"/>
				</div>
			</div>
			<!-- Photo modals -->
			<div v-for="photo in photos" :key="photo.id" :id="`photo-modal-${photo.id}`" class="modal fade" role="dialog" tabindex="-1">
				<div class="modal-dialog">
					<div class="modal-content">
						<div class="modal-body">
							<PhotoListItem :photo="photo" @error="(err) => this.errorMessage = err" @refresh="refresh"/>
						</div>
					</div>
				</div>
			</div>

			<ShowMore v-if="photosCursor && !loading" @loadMore="loadMore"/>
		</div>
	</PageSkeleton>
</template>

<style scoped>
.posts-grid {
	max-width: 480px;
	margin: 0 auto;
	display: grid;
	grid-template-columns: repeat(3, 1fr);
	gap: 10px;
	place-items: stretch;
}

.posts-grid .posts-grid-item {
	background-size: cover;
	background-position: center;
	background-repeat: no-repeat;
	width: 100%;
	height: 150px;
	cursor: pointer;
}
</style>
