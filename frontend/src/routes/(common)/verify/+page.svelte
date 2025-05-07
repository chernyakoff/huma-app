<script lang="ts">
	import { goto } from "$app/navigation";

	const email = localStorage.getItem("verify");

	if (!email) {
		goto("/login");
	}

	/* 
если токен есть
    если верифицирован, то показывем табличку мол успех и кнопку логин
    если нет, то отправляем на регистрацию
        в регистрации - если email существует, то ошибка а если существует но не вери


*/

	const timeout = 10;
	let disabled = $state(true);
	let countdown = $state(timeout);

	let interval: ReturnType<typeof setInterval>;

	function startCountdown() {
		interval = setInterval(() => {
			if (countdown > 0) {
				countdown -= 1;
			} else {
				disabled = false;
				clearInterval(interval);
			}
		}, 1000);
	}

	function handleClick() {
		if (!disabled) {
			alert("Кнопка нажата!");
			disabled = true;
			countdown = timeout;
			startCountdown();
		}
	}

	startCountdown();
</script>

<div class="hero flex flex-grow items-center justify-center">
	<div class="w-full max-w-sm shrink-0">
		<div class="card bg-base-100 shadow-sm">
			<div class="card-body">
				<h2 class="card-title">Registration Successful!</h2>
				<p>
					Thank you for registering. A confirmation email has been sent to your inbox. Please check
					your email to verify your account.
				</p>
				<div class="card-actions justify-end">
					<button class="btn btn-accent" onclick={handleClick} {disabled}>
						{disabled ? `Подождите ${countdown} секунд` : "Отправить повторно"}
					</button>
				</div>
			</div>
		</div>
	</div>
</div>

<style>
	/* button:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	} */
</style>
