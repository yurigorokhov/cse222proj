package com.ucsd.yurig.tcpbench;

import android.content.Context;
import android.content.pm.PackageInstaller;
import android.os.AsyncTask;
import android.os.PowerManager;
import android.os.StrictMode;
import android.support.v7.app.ActionBarActivity;
import android.os.Bundle;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.TextView;

import java.io.Console;
import java.io.DataInputStream;
import java.io.DataOutputStream;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.net.Socket;


public class MainActivity extends ActionBarActivity {

    class TcpConnector extends AsyncTask<String, String, String> {

        public Exception TcpEx = null;
        private DataOutputStream out = null;
        private String SessionToken;

        @Override
        protected String doInBackground(String... sessionToken) {
            SessionToken = sessionToken[0];
            try {
                OutputStream outToServer = socket.getOutputStream();
                out = new DataOutputStream(outToServer);
                InputStream inFromServer = socket.getInputStream();
                DataInputStream in = new DataInputStream(inFromServer);
                if(SessionToken == null) {
                    out.writeBytes("start\n");
                } else {
                    out.writeBytes("start " + SessionToken + "\n");
                }
                out.flush();
                String input;

                // read back session token
                SessionToken = in.readLine();
                while(!socket.isClosed() && (input = in.readLine()) != null) {
                    publishProgress(input);
                    out.writeBytes(SessionToken + "@" + input.trim() + "\n");
                    out.flush();
                }
            } catch(Exception ex) {
                TcpEx = ex;
            }
            return "Finished";
        }

        @Override
        protected void onPostExecute(String result) {
            if(TcpEx != null) {

                // keep going!
                status.setText(TcpEx.toString());
                Connect(hostname, port);
                connector.execute(SessionToken);
            }
        }

        @Override
        protected void onProgressUpdate(String... params) {
            status.setText(params[0]);
        }

        protected void finishUp() {
            try {
                out.writeBytes("close\n");
                out.flush();
                out.close();
                socket.close();
            } catch(Exception ex) {
                if(TcpEx != null) {
                    TcpEx = ex;
                }
            }
        }
    }

    private TextView status = null;
    public Socket socket = null;
    TcpConnector connector = null;
    String hostname = null;
    final int port = 11111;


    public void Connect(String hostname, int port) {
        try {
            if(socket != null && !socket.isClosed()) {
                while(!socket.isClosed())
                    socket.close();
            }
            socket = new Socket();
            connector = new TcpConnector();
            socket.connect(new InetSocketAddress(hostname, port));
            status.setText(String.format("Connected to %s", socket.getRemoteSocketAddress()));
        } catch(Exception e) {
            status.setText(e.toString());
        }
    }

    @Override
    protected void onCreate(Bundle savedInstanceState) {

        if (android.os.Build.VERSION.SDK_INT > 9) {
            StrictMode.ThreadPolicy policy = new StrictMode.ThreadPolicy.Builder().permitAll().build();
            StrictMode.setThreadPolicy(policy);
        }
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        final Button buttonConnect = (Button) findViewById(R.id.button_connect);
        final Button buttonStart = (Button) findViewById(R.id.button_start);
        final Button buttonStop = (Button) findViewById(R.id.button_stop);
        status = (TextView)findViewById(R.id.text_status);
        String token = null;

        // CONNECT
        buttonConnect.setOnClickListener(new View.OnClickListener() {
            public void onClick(View v) {


                hostname = ((EditText)findViewById(R.id.input_host)).getText().toString();
                status.setText(String.format("Connecting to %s", hostname));
                Connect(hostname, port);
            }
        });

        // START
        buttonStart.setOnClickListener(new View.OnClickListener() {
            public void onClick(View v) {
                if(socket == null || !socket.isConnected() || socket.isClosed() || connector == null || connector.isCancelled()) {
                    status.setText("Must connect first ... ");
                    return;
                }
                status.setText("Starting ... ");
                connector.execute((String)null);
            }
        });

        // STOP
        buttonStop.setOnClickListener(new View.OnClickListener() {
            public void onClick(View v) {
                try {
                    if(socket == null || !socket.isConnected() || socket.isClosed()) {
                        status.setText("Must connect first ... ");
                        return;
                    }
                    connector.cancel(true);
                    connector.finishUp();
                } catch(Exception e) {
                    status.setText(e.toString());
                }
                status.setText("Disconnected");
            }
        });
    }

    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        // Inflate the menu; this adds items to the action bar if it is present.
        getMenuInflater().inflate(R.menu.menu_main, menu);
        return true;
    }

    @Override
    public boolean onOptionsItemSelected(MenuItem item) {
        // Handle action bar item clicks here. The action bar will
        // automatically handle clicks on the Home/Up button, so long
        // as you specify a parent activity in AndroidManifest.xml.
        int id = item.getItemId();

        //noinspection SimplifiableIfStatement
        if (id == R.id.action_settings) {
            return true;
        }

        return super.onOptionsItemSelected(item);
    }
}
