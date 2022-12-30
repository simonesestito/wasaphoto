<script>
import {FollowService} from "../services";
import router from "../router";
import {toRefs} from "vue";

export default {
	name: "UserListItem",
	props: ['user'],
	emits: ['error'],
	data() {
		return {
			userData: this.user, // Use another internal variable for changes
			loading: false,
		};
	},
	setup(props) {
		return toRefs(props);
	},
	methods: {
		async follow() {
			await this.doFollow(true);
		},
		async unfollow() {
			await this.doFollow(false);
		},
		async doFollow(follow) {
			this.loading = true;
			this.$emit('error', '');
			try {
				if (follow) {
					await FollowService.followUser(this.user.id);
				} else {
					await FollowService.unfollowUser(this.user.id);
				}

				this.userData.following = follow;
			} catch (err) {
				this.$emit('error', err);
			} finally {
				this.loading = false;
			}
		},
		async goToUser() {
			await router.push(`/users/${this.userData.username}`);
		}
	}
}
</script>

<template>
	<div class="user-list-item p-4 mt-3"> <!-- @click="goToUser"> -->
		<div class="row">
			<div class="col col-lg-10 d-flex align-items-center" @click="goToUser">
				<p>
					<span><b>{{ userData.name }} {{ userData.surname }}</b></span>
					<br>
					<span>@{{ userData.username }}</span>
				</p>
			</div>
			<div class="col-md-auto d-flex align-items-center">
				<button @click="follow" v-if="!userData.following && !userData.banned" :disabled="loading" type="button"
						class="btn btn-outline-primary">Follow
				</button>
				<button @click="unfollow" v-if="userData.following && !userData.banned" :disabled="loading"
						type="button"
						class="btn btn-outline-secondary">Unfollow
				</button>
			</div>
		</div>
	</div>
</template>

<style scoped>
.user-list-item {
	cursor: pointer;
}

.user-list-item:hover {
	background-color: lightcyan;
}

.align-items-center p {
	margin: auto 0;
}
</style>

