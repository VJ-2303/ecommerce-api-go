// Register Page

document.addEventListener('DOMContentLoaded', () => {
    // Redirect if already authenticated
    if (isAuthenticated()) {
        window.location.href = '../index.html';
        return;
    }

    const registerForm = document.getElementById('registerForm');
    const passwordToggle = document.getElementById('passwordToggle');
    const confirmPasswordToggle = document.getElementById('confirmPasswordToggle');
    const passwordInput = document.getElementById('password');
    const confirmPasswordInput = document.getElementById('confirmPassword');

    // Password toggles
    if (passwordToggle && passwordInput) {
        passwordToggle.addEventListener('click', () => {
            const type = passwordInput.type === 'password' ? 'text' : 'password';
            passwordInput.type = type;
            passwordToggle.querySelector('.toggle-icon').textContent = 
                type === 'password' ? '👁️' : '🙈';
        });
    }

    if (confirmPasswordToggle && confirmPasswordInput) {
        confirmPasswordToggle.addEventListener('click', () => {
            const type = confirmPasswordInput.type === 'password' ? 'text' : 'password';
            confirmPasswordInput.type = type;
            confirmPasswordToggle.querySelector('.toggle-icon').textContent = 
                type === 'password' ? '👁️' : '🙈';
        });
    }

    // Form submission
    if (registerForm) {
        registerForm.addEventListener('submit', handleRegister);
    }
});

async function handleRegister(e) {
    e.preventDefault();
    clearAllErrors();

    const submitBtn = document.getElementById('submitBtn');
    const name = document.getElementById('name').value.trim();
    const phoneNumber = document.getElementById('phoneNumber').value.trim();
    const password = document.getElementById('password').value;
    const confirmPassword = document.getElementById('confirmPassword').value;

    // Client-side validation
    let hasError = false;

    if (name.length < 6) {
        showFormError('name', 'Name must be at least 6 characters');
        hasError = true;
    }

    if (!isValidPhone(phoneNumber)) {
        showFormError('phoneNumber', 'Please enter a valid 10-digit phone number');
        hasError = true;
    }

    if (password.length < 8) {
        showFormError('password', 'Password must be at least 8 characters');
        hasError = true;
    }

    if (password.length > 72) {
        showFormError('password', 'Password must be less than 72 characters');
        hasError = true;
    }

    if (password !== confirmPassword) {
        showFormError('confirmPassword', 'Passwords do not match');
        hasError = true;
    }

    if (hasError) return;

    setLoadingState(submitBtn, true);

    try {
        const data = await authApi.register({
            name: name,
            phone_number: phoneNumber,
            password: password
        });

        if (data.user) {
            showToast('Registration successful! Please login.', 'success');
            setTimeout(() => {
                window.location.href = 'login.html';
            }, 1500);
        }
    } catch (error) {
        console.error('Registration error:', error);
        handleApiError(error, {
            phone_number: 'phoneNumber'
        });
    } finally {
        setLoadingState(submitBtn, false);
    }
}
