from subprocess import check_output, call, Popen
from subprocess import call
from email.mime.text import MIMEText
import smtplib, string, time, thread
import datetime

def timeStamped(fname, fmt='%Y-%m-%d-%H-%M-%S_{fname}'):
    return datetime.datetime.now().strftime(fmt).format(fname=fname)

def sendEmail(toEmail, message, subject):
	try:
		sender = constants.EMAIL_USERNAME
		password = constants.EMAIL_PASSWORD

		msg = MIMEText(message, 'html')
		msg['Subject'] = subject
		msg['From'] = "signup@golinkstr.com"
		msg['To'] = toEmail

		smtpConn = smtplib.SMTP("smtp.1and1.com", 25)
		smtpConn.starttls()
		logger.debug("Attepting Login")
		smtpConn.login(sender, password)
		logger.debug("Sending email to " + toEmail)
		smtpConn.set_debuglevel(1)
		smtpConn.sendmail(sender, toEmail, msg.as_string())
		smtpConn.quit()
	except Exception, e:
		logger.error(e)
		return

def writeLog(message):
	with open('python_log.txt', 'a') as f:
		f.write(message)


while True:
	try:
		time.sleep(10)
		tmp = check_output(["ps", "-e"])
		if "links-" not in tmp:
			writeLog("Server is not running")
			Popen(["./links-as-a-service-redirectserver", "&"])
	except Exception, e:
		writeLog("Fuck, restart script failed")
