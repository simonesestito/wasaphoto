<template>
	<div class="row user-header" data-bs-dismiss="modal">
		<p @click="goToUser">{{ user.name }} {{ user.surname }} (@<u>{{ user.username }}</u>)</p>
	</div>
</template>

<script>
import router from "../router";

export default {
	name: "UserNameHeader",
	props: ['user'],
	methods: {
		goToUser() {
			router.push(`/users/${this.user.username}`);
		}
	}
}
</script>

<style scoped>
.user-header {
	height: 40px;
	font-size: 1.2rem;
}

.user-header * {
	cursor: pointer;
	width: auto;
}
</style>
