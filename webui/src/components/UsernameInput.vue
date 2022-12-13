<script>
export default {
	emits: [
		'submit',
	],
	props: [
		'loading',
		'submitText',
		'initialInput',
	],
	data() {
		return {
			userInput: this.initialInput || '',
			errorMessage: null,
		};
	},
	methods: {
		login(event) {
			event.preventDefault();
			if (this.validate())
				this.$emit('submit', this.userInput);
		},
		validate() {
			this.userInput = this.userInput.toLowerCase();

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
	<!-- Error -->
	<ErrorMsg v-if="errorMessage" :msg="errorMessage" />

	<form @submit="login">
		<div class="input-group">
			<div class="input-group-prepend">
				<span class="input-group-text" id="search-username-field">@</span>
			</div>
			<input type="text" class="form-control" placeholder="Username" aria-label="Username" v-model="userInput">
			<div class="input-group-append">
				<input class="btn btn-outline-primary" type="submit" :value="submitText" :disabled="loading">
			</div>
		</div>
	</form>
</template>
