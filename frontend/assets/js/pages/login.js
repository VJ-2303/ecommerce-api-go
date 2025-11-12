// Login Page

document.addEventListener('DOMContentLoaded', () => {
    // Redirect if already authenticated
    if (isAuthenticated()) {
        window.location.href = '../index.html';
        return;
    }

    const loginForm = document.getElementById('loginForm');
    const passwordToggle = document.getElementById('passwordToggle');
    const passwordInput = document.getElementById('password');

    // Password toggle
    if (passwordToggle && passwordInput) {
        passwordToggle.addEventListener('click', () => {
            const type = passwordInput.type === 'password' ? 'text' : 'password';
            passwordInput.type = type;
            passwordToggle.querySelector('.toggle-icon').textContent = 
                type === 'password' ? '👁️' : '🙈';
        });
    }

    // Form submission
    if (loginForm) {
        loginForm.addEventListener('submit', handleLogin);
    }
});

async function handleLogin(e) {
    e.preventDefault();
    clearAllErrors();

    const submitBtn = document.getElementById('submitBtn');
    const phoneNumber = document.getElementById('phoneNumber').value.trim();
    const password = document.getElementById('password').value;

    // Client-side validation
    let hasError = false;

    if (!isValidPhone(phoneNumber)) {
        showFormError('phoneNumber', 'Please enter a valid 10-digit phone number');
        hasError = true;
    }

    if (password.length < 8) {
        showFormError('password', 'Password must be at least 8 characters');
        hasError = true;
    }

    if (hasError) return;

    setLoadingState(submitBtn, true);

    try {
        const data = await authApi.login({
            phone_number: phoneNumber,
            password: password
        });

        if (data.auth_token) {
            // Parse token to get role
            const payload = parseJWT(data.auth_token.token);
            const role = payload?.role || 'user';
            
            saveAuthData(data.auth_token.token, role);
            showToast('Login successful!', 'success');

            setTimeout(() => {
                window.location.href = '../index.html';
            }, 1000);
        }
    } catch (error) {
        console.error('Login error:', error);
        handleApiError(error, {
            phone_number: 'phoneNumber'
        });
    } finally {
        setLoadingState(submitBtn, false);
    }
}
