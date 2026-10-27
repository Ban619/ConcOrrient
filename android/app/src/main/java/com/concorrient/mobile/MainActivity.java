package com.concorrient.mobile;

import android.app.Activity;
import android.app.AlertDialog;
import android.content.Intent;
import android.content.SharedPreferences;
import android.graphics.Color;
import android.net.Uri;
import android.os.Bundle;
import android.text.InputType;
import android.view.Gravity;
import android.view.ViewGroup;
import android.webkit.WebChromeClient;
import android.webkit.WebResourceRequest;
import android.webkit.WebSettings;
import android.webkit.WebView;
import android.webkit.WebViewClient;
import android.widget.Button;
import android.widget.EditText;
import android.widget.LinearLayout;
import android.widget.TextView;
import android.widget.Toast;

import java.net.URI;

public final class MainActivity extends Activity {
    private static final String PREFS_NAME = "concorrient";
    private static final String SERVER_KEY = "server_url";

    private WebView webView;
    private SharedPreferences preferences;
    private String serverUrl = "";

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        preferences = getSharedPreferences(PREFS_NAME, MODE_PRIVATE);
        serverUrl = preferences.getString(SERVER_KEY, "");
        buildScreen();
        if (serverUrl.isEmpty()) {
            showServerDialog(false);
        } else {
            loadWorkspace();
        }
    }

    private void buildScreen() {
        LinearLayout root = new LinearLayout(this);
        root.setOrientation(LinearLayout.VERTICAL);
        root.setBackgroundColor(Color.rgb(246, 247, 242));

        LinearLayout toolbar = new LinearLayout(this);
        toolbar.setGravity(Gravity.CENTER_VERTICAL);
        toolbar.setPadding(dp(16), dp(8), dp(12), dp(8));
        toolbar.setBackgroundColor(Color.rgb(39, 61, 54));

        TextView title = new TextView(this);
        title.setText("ConcOrrient");
        title.setTextColor(Color.WHITE);
        title.setTextSize(18);
        title.setTypeface(null, android.graphics.Typeface.BOLD);
        toolbar.addView(title, new LinearLayout.LayoutParams(0, dp(44), 1));

        Button connectionButton = new Button(this);
        connectionButton.setText("Server");
        connectionButton.setOnClickListener(view -> showServerDialog(true));
        toolbar.addView(connectionButton, new LinearLayout.LayoutParams(
                ViewGroup.LayoutParams.WRAP_CONTENT, dp(44)));

        webView = new WebView(this);
        WebSettings settings = webView.getSettings();
        settings.setJavaScriptEnabled(true);
        settings.setDomStorageEnabled(true);
        settings.setDatabaseEnabled(true);
        settings.setAllowFileAccess(false);
        settings.setAllowContentAccess(false);
        settings.setMixedContentMode(WebSettings.MIXED_CONTENT_NEVER_ALLOW);
        webView.setWebChromeClient(new WebChromeClient());
        webView.setWebViewClient(new WebViewClient() {
            @Override
            public boolean shouldOverrideUrlLoading(WebView view, WebResourceRequest request) {
                Uri target = request.getUrl();
                Uri workspace = Uri.parse(serverUrl);
                if (target.getHost() != null && target.getHost().equalsIgnoreCase(workspace.getHost())) {
                    return false;
                }
                try {
                    startActivity(new Intent(Intent.ACTION_VIEW, target));
                } catch (Exception ignored) {
                    Toast.makeText(MainActivity.this, "This link could not be opened.", Toast.LENGTH_SHORT).show();
                }
                return true;
            }

            @Override
            public void onReceivedError(WebView view, int errorCode, String description, String failingUrl) {
                if (view != null && failingUrl != null && failingUrl.equals(serverUrl + "/")) {
                    Toast.makeText(MainActivity.this, "Could not reach the workspace. Check the server and Wi-Fi address.", Toast.LENGTH_LONG).show();
                }
            }
        });

        root.addView(toolbar, new LinearLayout.LayoutParams(
                ViewGroup.LayoutParams.MATCH_PARENT, ViewGroup.LayoutParams.WRAP_CONTENT));
        root.addView(webView, new LinearLayout.LayoutParams(
                ViewGroup.LayoutParams.MATCH_PARENT, 0, 1));
        setContentView(root);
    }

    private void showServerDialog(boolean allowCancel) {
        EditText input = new EditText(this);
        input.setSingleLine(true);
        input.setInputType(InputType.TYPE_CLASS_TEXT | InputType.TYPE_TEXT_VARIATION_URI);
        input.setHint("http://192.168.1.20:8080");
        input.setText(serverUrl);
        input.setSelection(input.getText().length());

        LinearLayout content = new LinearLayout(this);
        content.setOrientation(LinearLayout.VERTICAL);
        content.setPadding(dp(24), dp(4), dp(24), 0);
        TextView instructions = new TextView(this);
        instructions.setText("Start ConcOrrient on your computer, then enter its Wi-Fi address. Both devices need to be on the same network.");
        instructions.setTextSize(14);
        content.addView(instructions);
        LinearLayout.LayoutParams inputParams = new LinearLayout.LayoutParams(
                ViewGroup.LayoutParams.MATCH_PARENT, ViewGroup.LayoutParams.WRAP_CONTENT);
        inputParams.topMargin = dp(12);
        content.addView(input, inputParams);

        AlertDialog.Builder builder = new AlertDialog.Builder(this)
                .setTitle("Connect to your workspace")
                .setView(content)
                .setPositiveButton("Connect", null);
        if (allowCancel) {
            builder.setNegativeButton("Cancel", (dialog, which) -> dialog.dismiss());
        }
        AlertDialog dialog = builder.create();
        dialog.setOnShowListener(ignored -> dialog.getButton(AlertDialog.BUTTON_POSITIVE).setOnClickListener(view -> {
            String value = input.getText().toString().trim().replaceAll("/+$", "");
            if (!isValidServerUrl(value)) {
                input.setError("Enter a complete http:// or https:// address.");
                return;
            }
            serverUrl = value;
            preferences.edit().putString(SERVER_KEY, serverUrl).apply();
            dialog.dismiss();
            loadWorkspace();
        }));
        dialog.show();
    }

    private boolean isValidServerUrl(String value) {
        try {
            URI uri = new URI(value);
            return ("http".equalsIgnoreCase(uri.getScheme()) || "https".equalsIgnoreCase(uri.getScheme()))
                    && uri.getHost() != null && uri.getUserInfo() == null;
        } catch (Exception ignored) {
            return false;
        }
    }

    private void loadWorkspace() {
        if (webView != null && !serverUrl.isEmpty()) {
            webView.loadUrl(serverUrl + "/");
        }
    }

    private int dp(int value) {
        return Math.round(value * getResources().getDisplayMetrics().density);
    }

    @Override
    public void onBackPressed() {
        if (webView != null && webView.canGoBack()) {
            webView.goBack();
        } else {
            super.onBackPressed();
        }
    }
}