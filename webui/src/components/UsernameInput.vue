<script>
export default {
	emits: [
		'submit',
		'error',
	],
	props: [
		'loading',
		'submitText',
	],
	data() {
		return {
			userInput: '',
			errorMessage: null,
		};
	},
	methods: {
		login(event) {
			event.preventDefault();
			if (this.validate())
				this.$emit('submit', this.userInput);
			else
				this.$emit('error', this.errorMessage);
		},
		validate() {
			if (this.userInput.length < 3)
				this.errorMessage = 'Username too short';
			else if (this.userInput.length > 16)
				this.errorMessage = 'Username too long';
			else if (!this.userInput.match(/^[a-z_0-9]+$/))
				this.errorMessage = 'Username not valid';
			else
				this.errorMessage = null;

			return this.errorMessage == null;
		},
	},
};
</script>

<template>
	<form @submit="login">
		<div class="input-group w-50">
			<div class="input-group-prepend">
				<span class="input-group-text" id="search-username-field">@</span>
			</div>
			<input type="text" class="form-control" placeholder="Username" aria-label="Username" v-model="userInput">
			<div class="input-group-append">
				<input class="btn btn-outline-primary" type="submit" :value="submitText">
			</div>
		</div>
	</form>
</template>
